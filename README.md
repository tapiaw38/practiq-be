# practiq-be

Backend de Practiq — plataforma educativa tipo Kumon construida en Go con arquitectura limpia.

## Requisitos

- Go 1.24+
- PostgreSQL 14+
- [golang-migrate CLI](https://github.com/golang-migrate/migrate) (para migraciones manuales)
- `auth-api-be` corriendo en el puerto 8082

## Configuración

Copia `.env.example` a `.env` y ajusta los valores:

```bash
cp .env.example .env
```

Variables principales:

| Variable     | Descripción                                    | Default          |
|-------------|------------------------------------------------|------------------|
| `PORT`       | Puerto HTTP del servicio                       | `8083`           |
| `DB_URL`     | DSN de PostgreSQL                              | ver `.env.example` |
| `JWT_SECRET` | Debe coincidir con el de `auth-api-be`         | —                |

## Base de datos

Levantar la base de datos con Docker:

```bash
docker-compose up -d practiq-postgres-db
```

O usar la instancia PostgreSQL existente ajustando `DB_URL` en `.env`.

## Migraciones

```bash
# Aplicar todas las migraciones
make migrate-up

# Revertir la última migración
make migrate-down

# Crear nueva migración
make migrate-create name=nombre_migracion
```

Las migraciones están en `migrations/`. La primera (`000001_init_schema`) crea todas las tablas y la segunda (`000002_seed_data`) inserta datos iniciales (estrategias Kumon, usuarios demo).

**Usuarios demo creados por el seed:**
- `teacher_demo` — rol `teacher`
- `student_demo` — rol `student`

> Estos usernames deben existir en `auth-api-be` para poder iniciar sesión.

## Ejecución

```bash
# Desarrollo (con live reload via air si está instalado)
make run

# Compilar y ejecutar el binario
make build
./build/practiq-be
```

El servidor queda disponible en `http://localhost:8083`.

## Docker Compose

Para levantar todo el stack (API + base de datos):

```bash
docker-compose up -d
```

## Arquitectura

```
practiq-be/
├── cmd/api/            # Punto de entrada (main.go)
├── internal/
│   ├── adapters/
│   │   ├── datasources/repositories/  # Implementaciones SQL (PostgreSQL)
│   │   └── web/                       # Handlers HTTP (Gin) + middlewares + rutas
│   ├── domain/                        # Entidades de negocio
│   ├── platform/
│   │   ├── appcontext/                # Contexto de la aplicación
│   │   ├── auth/                      # Validación JWT
│   │   ├── database/                  # Conexión a PostgreSQL
│   │   └── strategy/                  # Lógica de aprendizaje (Kumon)
│   └── usecases/                      # Casos de uso por dominio
├── migrations/                        # Archivos SQL de migración
└── docker-compose.yml
```

## Endpoints principales

Todos bajo el prefijo `/api` y requieren Bearer token (JWT emitido por `auth-api-be`).

| Método | Ruta | Descripción |
|--------|------|-------------|
| `POST` | `/api/profile` | Sincronizar perfil de usuario |
| `GET` | `/api/profile` | Obtener perfil propio |
| `GET/POST` | `/api/courses` | Listar / crear cursos |
| `POST` | `/api/courses/:id/enroll` | Inscribirse a un curso |
| `GET/POST` | `/api/courses/:id/topics` | Temas del curso |
| `GET/POST` | `/api/topics/:id/exercises` | Ejercicios del tema |
| `GET/POST` | `/api/courses/:id/practice-sheets` | Hojas de práctica |
| `POST` | `/api/practice-sheets/:id/submit` | Enviar respuestas |
| `GET` | `/api/students/me/progress` | Progreso del estudiante |
| `POST` | `/api/ai/help` | Asistente IA (mock) |
