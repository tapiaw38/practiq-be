DELETE FROM course_learning_strategies cls
USING learning_strategies ls
WHERE cls.strategy_id = ls.id
  AND ls.code = 'kumon'
  AND cls.is_default = TRUE
  AND cls.config = '{}'::jsonb;
