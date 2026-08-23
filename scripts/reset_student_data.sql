-- Borra la actividad de alumnos: intentos, entregas de cuaderno y el progreso
-- derivado de ellos. NO borra cursos, temas, ejercicios, hojas, cuadernos ni
-- perfiles — solo lo que los alumnos generaron al resolver.
--
-- PARA AMBIENTES SIN ALUMNOS REALES. Es destructivo e irreversible: no hay
-- forma de reconstruir un intento borrado.
--
-- Corré esto DESPUÉS de desplegar, para que todo lo que quede lo haya escrito
-- el flujo actualizado:
--
--   psql "$DATABASE_URL" -f scripts/reset_student_data.sql
--
-- Reemplaza a la migración que limpiaba `needs_teacher_review` en prácticas:
-- resuelve también las pruebas de nivel que el `requireTeacher` viejo dejó
-- marcadas con score 0 pese a tener veredicto de la IA, y no deja un UPDATE
-- irreversible en el historial de migraciones.
BEGIN;

-- Los intentos y su canvas (student_work_canvas cae por ON DELETE CASCADE).
DELETE FROM student_attempts;

-- El candado de "esta prueba de nivel ya fue entregada". Sin borrarlo, los
-- alumnos quedan sin poder reenviar ninguna prueba aunque no exista el intento:
-- el submit responde "this level test was already submitted".
DELETE FROM level_test_submissions;

-- Entregas de cuaderno.
DELETE FROM notebook_submissions;

-- Progreso derivado de todo lo anterior: se recalcula solo al volver a resolver.
DELETE FROM student_topic_progress;
DELETE FROM student_course_progress;

-- Trabajos de envío asincrónico, en vuelo o terminados.
DELETE FROM submit_jobs;

-- Opcional: pedidos de ayuda a la IA y notificaciones ya emitidas. Sacá estas
-- dos líneas si querés conservar ese historial.
DELETE FROM ai_help_requests;
DELETE FROM notifications;

COMMIT;
