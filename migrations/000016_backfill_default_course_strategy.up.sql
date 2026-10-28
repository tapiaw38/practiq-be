INSERT INTO course_learning_strategies (course_id, strategy_id, is_default, config)
SELECT c.id, ls.id, TRUE, '{}'::jsonb
FROM courses c
JOIN learning_strategies ls ON ls.code = 'kumon' AND ls.status = 'active'
WHERE NOT EXISTS (
    SELECT 1
    FROM course_learning_strategies cls
    WHERE cls.course_id = c.id
);
