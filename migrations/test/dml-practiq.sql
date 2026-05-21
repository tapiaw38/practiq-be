INSERT INTO learning_strategies (name, code, description, status)
VALUES (
    'Kumon Inspired',
    'kumon',
    'Estrategia basada en practica diaria, repeticion, dominio progresivo y avance por niveles.',
    'active'
) ON CONFLICT (code) DO NOTHING;

INSERT INTO user_profiles (id, name, email, profile_type)
VALUES
    ('teacher_demo', 'Profesor Demo', 'teacher@practiq.com', 'teacher'),
    ('student_demo', 'Walter Tapia', 'student@practiq.com', 'student')
ON CONFLICT (id) DO NOTHING;

INSERT INTO courses (id, teacher_id, title, description, level, subject)
VALUES (
    '20000000-0000-0000-0000-000000000001',
    'teacher_demo',
    'Matematica base',
    'Curso demo para practicar fracciones, sumas y progresion guiada.',
    'Primaria',
    'Matematica'
) ON CONFLICT (id) DO NOTHING;

INSERT INTO enrollments (id, course_id, student_id, status)
VALUES (
    '30000000-0000-0000-0000-000000000001',
    '20000000-0000-0000-0000-000000000001',
    'student_demo',
    'active'
) ON CONFLICT (course_id, student_id) DO NOTHING;

INSERT INTO course_learning_strategies (id, course_id, strategy_id, is_default, config)
VALUES (
    '40000000-0000-0000-0000-000000000001',
    '20000000-0000-0000-0000-000000000001',
    (SELECT id FROM learning_strategies WHERE code = 'kumon'),
    TRUE,
    '{"daily_goal": 10}'
) ON CONFLICT (course_id, strategy_id) DO NOTHING;

INSERT INTO topics (id, course_id, title, description, order_index)
VALUES
    ('50000000-0000-0000-0000-000000000001', '20000000-0000-0000-0000-000000000001', 'Fracciones', 'Introduccion a sumas y restas de fracciones.', 1),
    ('50000000-0000-0000-0000-000000000002', '20000000-0000-0000-0000-000000000001', 'Numeros decimales', 'Practica guiada con equivalencias y orden.', 2)
ON CONFLICT (id) DO NOTHING;

INSERT INTO exercises (id, topic_id, type, question, correct_answer, explanation, difficulty, metadata)
VALUES
    ('60000000-0000-0000-0000-000000000001', '50000000-0000-0000-0000-000000000001', 'equation', '3/4 + 1/8', '7/8', 'Convierte 3/4 a octavos y luego suma.', 3, '{}'),
    ('60000000-0000-0000-0000-000000000002', '50000000-0000-0000-0000-000000000001', 'equation', '2/3 + 1/6', '5/6', 'Usa denominador comun 6.', 4, '{}'),
    ('60000000-0000-0000-0000-000000000003', '50000000-0000-0000-0000-000000000001', 'equation', '7/10 - 3/5', '1/10', 'Convierte 3/5 en 6/10.', 4, '{}'),
    ('60000000-0000-0000-0000-000000000004', '50000000-0000-0000-0000-000000000001', 'equation', '1/2 + 1/4 + 1/8', '7/8', 'Suma fracciones equivalentes con denominador 8.', 5, '{}'),
    ('60000000-0000-0000-0000-000000000005', '50000000-0000-0000-0000-000000000001', 'equation', '5/6 - 1/3', '1/2', '1/3 equivale a 2/6.', 5, '{}'),
    ('60000000-0000-0000-0000-000000000006', '50000000-0000-0000-0000-000000000002', 'equation', '0.5 + 0.25', '0.75', 'Suma las partes decimales.', 2, '{}'),
    ('60000000-0000-0000-0000-000000000007', '50000000-0000-0000-0000-000000000002', 'equation', '1.2 - 0.4', '0.8', 'Alinea las cifras decimales.', 2, '{}')
ON CONFLICT (id) DO NOTHING;

INSERT INTO practice_sheets (id, course_id, topic_id, strategy_id, title, level, sheet_type, test_style, created_by)
VALUES
    ('70000000-0000-0000-0000-000000000001', '20000000-0000-0000-0000-000000000001', '50000000-0000-0000-0000-000000000001', (SELECT id FROM learning_strategies WHERE code = 'kumon'), 'Practica guiada de fracciones', 3, 'practice', 'keyboard', 'teacher'),
    ('70000000-0000-0000-0000-000000000002', '20000000-0000-0000-0000-000000000001', '50000000-0000-0000-0000-000000000002', (SELECT id FROM learning_strategies WHERE code = 'kumon'), 'Practica inicial de decimales', 2, 'practice', 'keyboard', 'teacher'),
    ('70000000-0000-0000-0000-000000000003', '20000000-0000-0000-0000-000000000001', '50000000-0000-0000-0000-000000000001', (SELECT id FROM learning_strategies WHERE code = 'kumon'), 'Prueba de nivel 3 — Fracciones (hoja)', 3, 'level_test', 'canvas', 'teacher')
ON CONFLICT (id) DO NOTHING;

INSERT INTO practice_sheet_exercises (id, practice_sheet_id, exercise_id, order_index)
VALUES
    ('80000000-0000-0000-0000-000000000008', '70000000-0000-0000-0000-000000000003', '60000000-0000-0000-0000-000000000001', 1),
    ('80000000-0000-0000-0000-000000000009', '70000000-0000-0000-0000-000000000003', '60000000-0000-0000-0000-000000000002', 2),
    ('80000000-0000-0000-0000-000000000010', '70000000-0000-0000-0000-000000000003', '60000000-0000-0000-0000-000000000003', 3),
    ('80000000-0000-0000-0000-000000000011', '70000000-0000-0000-0000-000000000003', '60000000-0000-0000-0000-000000000004', 4),
    ('80000000-0000-0000-0000-000000000012', '70000000-0000-0000-0000-000000000003', '60000000-0000-0000-0000-000000000005', 5)
ON CONFLICT (practice_sheet_id, exercise_id) DO NOTHING;

INSERT INTO practice_sheet_exercises (id, practice_sheet_id, exercise_id, order_index)
VALUES
    ('80000000-0000-0000-0000-000000000001', '70000000-0000-0000-0000-000000000001', '60000000-0000-0000-0000-000000000001', 1),
    ('80000000-0000-0000-0000-000000000002', '70000000-0000-0000-0000-000000000001', '60000000-0000-0000-0000-000000000002', 2),
    ('80000000-0000-0000-0000-000000000003', '70000000-0000-0000-0000-000000000001', '60000000-0000-0000-0000-000000000003', 3),
    ('80000000-0000-0000-0000-000000000004', '70000000-0000-0000-0000-000000000001', '60000000-0000-0000-0000-000000000004', 4),
    ('80000000-0000-0000-0000-000000000005', '70000000-0000-0000-0000-000000000001', '60000000-0000-0000-0000-000000000005', 5),
    ('80000000-0000-0000-0000-000000000006', '70000000-0000-0000-0000-000000000002', '60000000-0000-0000-0000-000000000006', 1),
    ('80000000-0000-0000-0000-000000000007', '70000000-0000-0000-0000-000000000002', '60000000-0000-0000-0000-000000000007', 2)
ON CONFLICT (practice_sheet_id, exercise_id) DO NOTHING;

-- ── Notebooks ───────────────────────────────────────────────────────────────

INSERT INTO notebooks (id, course_id, teacher_id, title, description)
VALUES
    (
        'a0000000-0000-0000-0000-000000000001',
        '20000000-0000-0000-0000-000000000001',
        'teacher_demo',
        'Cuaderno de Fracciones',
        'Practica diaria de sumas y restas de fracciones. Completa cada pagina con tu respuesta.'
    ),
    (
        'a0000000-0000-0000-0000-000000000002',
        '20000000-0000-0000-0000-000000000001',
        'teacher_demo',
        'Cuaderno de Decimales',
        'Ejercicios de numeros decimales para practicar en casa.'
    )
ON CONFLICT (id) DO NOTHING;

INSERT INTO notebook_pages (id, notebook_id, page_number, title, content_type, content_data, instructions)
VALUES
    (
        'b0000000-0000-0000-0000-000000000001',
        'a0000000-0000-0000-0000-000000000001',
        1,
        'Suma de fracciones con igual denominador',
        'text',
        'Resuelve las siguientes sumas:

1) 1/5 + 2/5 = ___
2) 3/8 + 4/8 = ___
3) 2/7 + 3/7 = ___
4) 5/9 + 1/9 = ___
5) 4/11 + 6/11 = ___',
        'Escribe el resultado de cada suma en el espacio en blanco.'
    ),
    (
        'b0000000-0000-0000-0000-000000000002',
        'a0000000-0000-0000-0000-000000000001',
        2,
        'Suma con distinto denominador',
        'text',
        'Encuentra el denominador comun y resuelve:

1) 1/2 + 1/4 = ___
2) 1/3 + 1/6 = ___
3) 2/5 + 1/10 = ___
4) 3/4 + 1/8 = ___
5) 1/2 + 1/3 = ___',
        'Recuerda encontrar el minimo comun multiplo antes de sumar.'
    ),
    (
        'b0000000-0000-0000-0000-000000000003',
        'a0000000-0000-0000-0000-000000000001',
        3,
        'Resta de fracciones',
        'text',
        'Calcula las siguientes restas:

1) 3/4 - 1/4 = ___
2) 7/8 - 3/8 = ___
3) 5/6 - 1/3 = ___
4) 2/3 - 1/6 = ___
5) 9/10 - 2/5 = ___',
        'Convierte a un denominador comun cuando sea necesario.'
    ),
    (
        'b0000000-0000-0000-0000-000000000004',
        'a0000000-0000-0000-0000-000000000002',
        1,
        'Suma de decimales',
        'text',
        'Resuelve alineando los puntos decimales:

1) 0.5 + 0.3 = ___
2) 1.2 + 0.8 = ___
3) 3.14 + 1.06 = ___
4) 0.75 + 0.25 = ___
5) 2.5 + 1.35 = ___',
        'Alinea siempre la coma decimal antes de sumar.'
    ),
    (
        'b0000000-0000-0000-0000-000000000005',
        'a0000000-0000-0000-0000-000000000002',
        2,
        'Resta de decimales',
        'text',
        'Calcula:

1) 1.0 - 0.4 = ___
2) 2.5 - 1.2 = ___
3) 5.0 - 3.75 = ___
4) 10.0 - 4.6 = ___
5) 3.33 - 1.11 = ___',
        'Alinea siempre la coma decimal antes de restar.'
    )
ON CONFLICT (id) DO NOTHING;

INSERT INTO student_topic_progress (
    id,
    student_id,
    topic_id,
    strategy_id,
    mastery_score,
    current_level,
    total_attempts,
    correct_attempts,
    streak_days,
    last_practiced_at,
    updated_at
)
VALUES
    (
        '90000000-0000-0000-0000-000000000001',
        'student_demo',
        '50000000-0000-0000-0000-000000000001',
        (SELECT id FROM learning_strategies WHERE code = 'kumon'),
        68,
        3,
        10,
        7,
        12,
        NOW() - INTERVAL '1 day',
        NOW()
    ),
    (
        '90000000-0000-0000-0000-000000000002',
        'student_demo',
        '50000000-0000-0000-0000-000000000002',
        (SELECT id FROM learning_strategies WHERE code = 'kumon'),
        42,
        2,
        6,
        4,
        5,
        NOW() - INTERVAL '2 day',
        NOW()
    )
ON CONFLICT (student_id, topic_id, strategy_id) DO UPDATE SET
    mastery_score = EXCLUDED.mastery_score,
    current_level = EXCLUDED.current_level,
    total_attempts = EXCLUDED.total_attempts,
    correct_attempts = EXCLUDED.correct_attempts,
    streak_days = EXCLUDED.streak_days,
    last_practiced_at = EXCLUDED.last_practiced_at,
    updated_at = EXCLUDED.updated_at;
