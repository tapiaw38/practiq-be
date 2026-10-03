ALTER TABLE ai_help_requests
  DROP CONSTRAINT IF EXISTS ai_help_requests_help_type_check;

ALTER TABLE ai_help_requests
  ADD CONSTRAINT ai_help_requests_help_type_check
  CHECK (help_type IN ('hint', 'explanation', 'similar_example'));
