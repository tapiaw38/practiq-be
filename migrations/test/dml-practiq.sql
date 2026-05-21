INSERT INTO learning_strategies (name, code, description, status)
VALUES ('Kumon Inspired', 'kumon', 'Estrategia basada en practica diaria, repeticion y dominio progresivo.', 'active')
ON CONFLICT (code) DO NOTHING;

INSERT INTO user_profiles (id, name, email, profile_type)
VALUES
    ('teacher_demo',            'Profesor Demo', 'teacher@practiq.com', 'teacher'),
    ('walter.tapia.srmeeypfkf', 'Walter Tapia',  'tapiaw38@gmail.com',  'student')
ON CONFLICT (id) DO NOTHING;

-- ── Curso ────────────────────────────────────────────────────────────────────

INSERT INTO courses (id, teacher_id, title, description, level, subject)
VALUES (
    '20000000-0000-0000-0000-000000000001',
    'teacher_demo',
    'Matematica base',
    'Curso de matematica con progresion por niveles. Fracciones, decimales y operaciones.',
    'Primaria',
    'Matematica'
) ON CONFLICT (id) DO NOTHING;

INSERT INTO enrollments (id, course_id, student_id, status)
VALUES ('30000000-0000-0000-0000-000000000001', '20000000-0000-0000-0000-000000000001', 'walter.tapia.srmeeypfkf', 'active')
ON CONFLICT (course_id, student_id) DO NOTHING;

INSERT INTO course_learning_strategies (id, course_id, strategy_id, is_default, config)
VALUES (
    '40000000-0000-0000-0000-000000000001',
    '20000000-0000-0000-0000-000000000001',
    (SELECT id FROM learning_strategies WHERE code = 'kumon'),
    TRUE,
    '{"daily_goal": 10}'
) ON CONFLICT (course_id, strategy_id) DO NOTHING;

-- ── Nivel 1: Fracciones basicas ───────────────────────────────────────────

INSERT INTO topics (id, course_id, title, description, order_index)
VALUES ('50000000-0000-0000-0000-000000000001', '20000000-0000-0000-0000-000000000001', 'Fracciones basicas', 'Introduccion a fracciones simples con igual denominador.', 1)
ON CONFLICT (id) DO NOTHING;

INSERT INTO exercises (id, topic_id, type, question, correct_answer, explanation, difficulty, metadata)
VALUES
    ('60000000-0000-0000-0000-000000000001', '50000000-0000-0000-0000-000000000001', 'equation', '1/4 + 2/4',  '3/4',  'Mismos denominadores: suma los numeradores.', 1, '{}'),
    ('60000000-0000-0000-0000-000000000002', '50000000-0000-0000-0000-000000000001', 'equation', '3/5 + 1/5',  '4/5',  'Mismos denominadores: suma los numeradores.', 1, '{}'),
    ('60000000-0000-0000-0000-000000000003', '50000000-0000-0000-0000-000000000001', 'equation', '5/8 - 2/8',  '3/8',  'Mismos denominadores: resta los numeradores.', 1, '{}'),
    ('60000000-0000-0000-0000-000000000004', '50000000-0000-0000-0000-000000000001', 'equation', '4/6 - 1/6',  '3/6',  'Mismos denominadores: resta los numeradores.', 1, '{}'),
    ('60000000-0000-0000-0000-000000000005', '50000000-0000-0000-0000-000000000001', 'equation', '2/7 + 3/7',  '5/7',  'Mismos denominadores: suma los numeradores.', 1, '{}'),
    ('60000000-0000-0000-0000-000000000006', '50000000-0000-0000-0000-000000000001', 'equation', '6/9 - 2/9',  '4/9',  'Mismos denominadores: resta los numeradores.', 2, '{}'),
    ('60000000-0000-0000-0000-000000000007', '50000000-0000-0000-0000-000000000001', 'equation', '1/3 + 1/3',  '2/3',  'Mismos denominadores: suma los numeradores.', 1, '{}'),
    ('60000000-0000-0000-0000-000000000008', '50000000-0000-0000-0000-000000000001', 'equation', '7/10 - 3/10','4/10', 'Mismos denominadores: resta los numeradores.', 2, '{}')
ON CONFLICT (id) DO NOTHING;

-- Practica nivel 1
INSERT INTO practice_sheets (id, course_id, topic_id, strategy_id, title, level, sheet_type, test_style, created_by)
VALUES (
    '70000000-0000-0000-0000-000000000001',
    '20000000-0000-0000-0000-000000000001',
    '50000000-0000-0000-0000-000000000001',
    (SELECT id FROM learning_strategies WHERE code = 'kumon'),
    'Practica de fracciones - Nivel 1',
    1, 'practice', 'keyboard', 'teacher'
) ON CONFLICT (id) DO NOTHING;

INSERT INTO practice_sheet_exercises (practice_sheet_id, exercise_id, order_index)
VALUES
    ('70000000-0000-0000-0000-000000000001', '60000000-0000-0000-0000-000000000001', 1),
    ('70000000-0000-0000-0000-000000000001', '60000000-0000-0000-0000-000000000002', 2),
    ('70000000-0000-0000-0000-000000000001', '60000000-0000-0000-0000-000000000003', 3),
    ('70000000-0000-0000-0000-000000000001', '60000000-0000-0000-0000-000000000004', 4),
    ('70000000-0000-0000-0000-000000000001', '60000000-0000-0000-0000-000000000005', 5)
ON CONFLICT (practice_sheet_id, exercise_id) DO NOTHING;

-- Prueba de nivel 1
INSERT INTO practice_sheets (id, course_id, topic_id, strategy_id, title, level, sheet_type, test_style, created_by)
VALUES (
    '70000000-0000-0000-0000-000000000002',
    '20000000-0000-0000-0000-000000000001',
    '50000000-0000-0000-0000-000000000001',
    (SELECT id FROM learning_strategies WHERE code = 'kumon'),
    'Prueba de Nivel 1',
    1, 'level_test', 'keyboard', 'teacher'
) ON CONFLICT (id) DO NOTHING;

INSERT INTO practice_sheet_exercises (practice_sheet_id, exercise_id, order_index)
VALUES
    ('70000000-0000-0000-0000-000000000002', '60000000-0000-0000-0000-000000000001', 1),
    ('70000000-0000-0000-0000-000000000002', '60000000-0000-0000-0000-000000000002', 2),
    ('70000000-0000-0000-0000-000000000002', '60000000-0000-0000-0000-000000000003', 3),
    ('70000000-0000-0000-0000-000000000002', '60000000-0000-0000-0000-000000000004', 4),
    ('70000000-0000-0000-0000-000000000002', '60000000-0000-0000-0000-000000000005', 5),
    ('70000000-0000-0000-0000-000000000002', '60000000-0000-0000-0000-000000000006', 6),
    ('70000000-0000-0000-0000-000000000002', '60000000-0000-0000-0000-000000000007', 7),
    ('70000000-0000-0000-0000-000000000002', '60000000-0000-0000-0000-000000000008', 8)
ON CONFLICT (practice_sheet_id, exercise_id) DO NOTHING;

-- Cuaderno nivel 1
INSERT INTO notebooks (id, course_id, teacher_id, title, description, level)
VALUES (
    'a0000000-0000-0000-0000-000000000001',
    '20000000-0000-0000-0000-000000000001',
    'teacher_demo',
    'Cuaderno de fracciones - Nivel 1',
    'Resuelve cada pagina escribiendo el procedimiento a mano.',
    1
) ON CONFLICT (id) DO NOTHING;

INSERT INTO notebook_pages (notebook_id, page_number, title, content_type, content_data, instructions)
VALUES
    ('a0000000-0000-0000-0000-000000000001', 1, 'Suma de fracciones iguales',
     'text',
     '1) 1/5 + 2/5 = ___
2) 3/8 + 4/8 = ___
3) 2/7 + 3/7 = ___
4) 1/4 + 2/4 = ___
5) 4/9 + 3/9 = ___',
     'Suma los numeradores y conserva el denominador.'),
    ('a0000000-0000-0000-0000-000000000001', 2, 'Resta de fracciones iguales',
     'text',
     '1) 5/8 - 2/8 = ___
2) 4/6 - 1/6 = ___
3) 7/9 - 3/9 = ___
4) 6/7 - 2/7 = ___
5) 9/10 - 4/10 = ___',
     'Resta los numeradores y conserva el denominador.')
ON CONFLICT DO NOTHING;

-- ── Progreso del alumno (nivel 1) ────────────────────────────────────────────

INSERT INTO student_course_progress (student_id, course_id, current_level)
VALUES ('walter.tapia.srmeeypfkf', '20000000-0000-0000-0000-000000000001', 1)
ON CONFLICT (student_id, course_id) DO UPDATE SET current_level = 1;
