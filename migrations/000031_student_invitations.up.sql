-- Código de invitación que el docente dicta a sus alumnos para vincularse sin
-- pasar por el administrador, único que hasta ahora podía crear la asignación.
CREATE TABLE IF NOT EXISTS student_invitations (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code        VARCHAR(16) NOT NULL UNIQUE,
    teacher_id  VARCHAR(255) NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    -- Reservadas: hoy el canje solo crea el vínculo docente-alumno. Están acá
    -- para que sumar matrícula o cupo después no necesite otra migración.
    course_id   UUID REFERENCES courses(id) ON DELETE SET NULL,
    grade_id    UUID REFERENCES grades(id) ON DELETE SET NULL,
    max_uses    INT,
    uses        INT NOT NULL DEFAULT 0,
    expires_at  TIMESTAMP,
    revoked_at  TIMESTAMP,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Un solo código vigente por docente. En el índice y no en el código de la
-- aplicación: dos pedidos simultáneos de regeneración no pueden dejar dos
-- códigos activos.
CREATE UNIQUE INDEX IF NOT EXISTS idx_student_invitations_active_per_teacher
    ON student_invitations (teacher_id) WHERE revoked_at IS NULL;

-- Quién entró con qué código. Hace el canje idempotente: repetirlo no vuelve a
-- sumar un uso.
CREATE TABLE IF NOT EXISTS invitation_redemptions (
    invitation_id UUID NOT NULL REFERENCES student_invitations(id) ON DELETE CASCADE,
    student_id    VARCHAR(255) NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    created_at    TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (invitation_id, student_id)
);

CREATE INDEX IF NOT EXISTS idx_invitation_redemptions_student
    ON invitation_redemptions (student_id);
