# practiq-be — usecase and handler structure

Canonical reference: `internal/usecases/course/create.go`.

## 1. One file per public use case

The file is named after the operation, not the domain:

```sh
internal/usecases/<domain>/
  create.go
  get.go
  list.go
  update.go
  delete.go
  archive.go
```

Never group operations in `manage.go`, `handlers.go`, `repository.go`, or any
file named after the domain instead of the operation. Every operation gets its
own file: `create.go`, `update.go`, `close.go`, `add_member.go`, and so on.
Umbrella files such as `write.go`, `lifecycle.go` or `members.go` do not meet
this rule either.

## 2. What goes in an operation file

```go
type (
    CreateUsecase interface {
        Execute(ctx context.Context, teacherID, schoolID string, isSuperAdmin bool, in CreateInput) (*CreateOutput, apperrors.ApplicationError)
    }

    createUsecase struct {
        contextFactory appcontext.Factory
    }

    CreateInput struct {
        GradeID string `json:"grade_id"`
        Title   string `json:"title"`
    }

    CreateOutput struct {
        Data CourseData `json:"data"`
    }
)

func NewCreateUsecase(contextFactory appcontext.Factory) CreateUsecase { ... }

func (u *createUsecase) Execute(...) (*CreateOutput, apperrors.ApplicationError) { ... }
```

- **One** interface per file, with a **single** public method: `Execute`.
- Private struct with the same name, lowercased.
- Constructor `New<Operation>Usecase`.
- `Input` and `Output` belong to the operation. Shared `Data` types live in
  `output_types.go`.

An interface with two or more methods means a file is missing.

## 3. Helpers

- A helper used by **one** operation lives in that operation's file.
- A helper shared by several goes in a file with an explicit name:

```sh
access.go            permissions and visibility
validation.go        input rules
output_types.go      shared response structs
scope.go             plan and scope resolution
```

A helper file is fine; a file holding several operations is not.

## 4. What a refactor must not change

Moving code is not changing behaviour. These stay identical:

- HTTP routes and their verbs.
- The shape of request and response JSON, field by field.
- HTTP status codes and internal error codes.
- Authorisation, validation, transactions, and the order of operations.

If none of that changes, **do not touch the frontend**.

## 5. After each domain

```bash
gofmt -w ./internal/...
go test ./internal/usecases/<domain>/...
go test ./...
git diff --check
```

## 6. Comments

This repo carries no explanatory comments in the code. The reasoning goes in
the commit message.

## 7. Errors

Error details are declared in `internal/platform/errors/mappings/domain.go`,
never inline in a usecase:

```go
// mappings/domain.go
GillieSettingsGetError = ErrorDetails{
    InternalCode: "gillie:read-error",
    StatusCode:   http.StatusInternalServerError,
    Message:      "could not read the assistant settings",
}

// usecase
return nil, apperrors.NewApplicationError(mappings.GillieSettingsGetError, err)
```

Use `http.StatusX` constants, not bare numbers. When the message has to be
dynamic, copy the mapping and override only the message, so the code and the
status still come from one place:

```go
details := mappings.GillieSettingsBadURLError
details.Message = err.Error()
return nil, apperrors.NewApplicationError(details, err)
```

## 8. Inputs live in the usecase

Reference: `auth-api-be`, `internal/usecases/user/change_password.go`.

The `Input` is the use case's contract and lives **in the usecase**, with
`json` tags. The handler binds straight into it.

```go
// handlers/course/create.go
var input ucCourse.CreateInput
if err := c.ShouldBindJSON(&input); err != nil { ... }

output, appErr := uc.Execute(c, middlewares.GetUserID(c), c.GetHeader("X-School-ID"), middlewares.IsSuperAdmin(c), input)
```

Rules:

- **The Input is the body and nothing else.** Nothing from the token, the path,
  a query param or a header belongs inside it.
- **Identity goes as parameters**, before the Input: `teacherID`, `schoolID`,
  `isSuperAdmin`, path ids. Never as struct fields.
- **There is no `CreateCommand` or `CreateParams`.** `json` tags do not tie an
  Input to HTTP: a job or a CLI builds the same struct and the tags are ignored.
- **No `input_types.go` in handlers.** Duplicating the struct to copy it field
  by field is work nobody reads.

What this rule prevents, which already happened in this repo:

```go
Execute(context.Context, bool, string, string, bool, string, string, string, string, string, string)
```

Eleven arguments, seven `string` in a row. Swapping two of them still compiles
and breaks in production. If a signature grows past roughly four parameters
plus the Input, it means body fields are loose: they belong in the Input.

## 9. HTTP handlers

One endpoint per file, named after the operation:

```sh
internal/adapters/web/handlers/<domain>/
  create.go
  get.go
  list.go
  update.go
  delete.go
```

- Each file exports only `New<Operation>Handler`.
- The handler does one job: bind the body into the usecase's Input, take
  identity from the middleware, call `Execute`, and respond.
- Business logic, authorisation and domain validation live in the usecase.
- Response types of its own only when they differ from the usecase output;
  otherwise return the output directly.
- Never create `handlers.go`, `handler.go`, or files holding several endpoints.

A handler refactor must not change routes, verbs, required headers, JSON,
status codes or error codes.

## 10. Repositories

Canonical reference: `internal/adapters/datasources/repositories/course/`.

One package per table or aggregate, one file per query, named after the
operation:

```sh
internal/adapters/datasources/repositories/<domain>/
  repository.go        interface, struct, constructor, filter options
  create.go
  get.go
  list.go
  update.go
  delete.go
```

`repository.go` holds only the seams, never SQL:

```go
type ListFilterOptions struct {
    TeacherID string
    SchoolID  string
}

type Repository interface {
    Create(context.Context, domain.Course) (string, error)
    Get(context.Context, string) (*domain.Course, error)
    List(context.Context, ListFilterOptions) ([]domain.Course, error)
}

type repository struct {
    db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
    return &repository{db: db}
}
```

### Domain types

An entity read from or written to the database belongs in `internal/domain`,
never in a repository package:

```go
// internal/domain/site_contact.go
type SiteContact struct {
    Email    string
    Phone    string
    WhatsApp string
}
```

The repository imports that domain type in its interface and operation files:

```go
type Repository interface {
    Get(context.Context) (domain.SiteContact, error)
    Save(context.Context, domain.SiteContact) (domain.SiteContact, error)
}
```

- Repository packages may define only repository seams: `Repository`, private
  `repository`, constructors, and query-specific filter/options types.
- Never define `Contact`, `Course`, `Student`, or another persisted business
  entity in `internal/adapters/datasources/repositories/**`.
- `json` tags do not belong on domain entities for repository convenience.
- A type is a repository filter only when it describes query mechanics, such
  as `ListFilterOptions`; otherwise it belongs in domain.

Each operation file holds one method with its query:

```go
func (r *repository) Create(ctx context.Context, c domain.Course) (string, error) {
    query := `
        INSERT INTO courses (teacher_id, title)
        VALUES ($1, $2)
        RETURNING id
    `
    var id string
    err := r.db.QueryRowContext(ctx, query, c.TeacherID, c.Title).Scan(&id)
    return id, err
}
```

Rules:

- **SQL lives here and nowhere else.** A usecase never builds a query.
- **Return domain types**, never rows or driver types. A row that does not exist
  is `nil, nil`, not an error — "no subscription" is a normal state.
- **Return the raw `error`.** `apperrors.ApplicationError` belongs to the
  usecase, which decides what a failed read means for the caller.
- **Filters are named structs** (`ListFilterOptions`), not a string of
  positional arguments — the same reason use case signatures avoid them.
- **Context is always the first argument** and reaches the driver through
  `QueryRowContext` / `QueryContext` / `ExecContext`, never `Query` or `Exec`.
- Query-only helpers (scanners, nullable wrappers) go in their own file, such
  as `nullable_uuid.go`.
- A new repository is registered in `repositories.go`, which is the only file
  that knows the whole set.

Small repositories with two or three methods may keep them in `repository.go`;
`site_contact` and `gillie_settings` do. Once a package grows past that, split
it by operation.
