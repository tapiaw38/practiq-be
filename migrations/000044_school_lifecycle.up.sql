-- Closing a school must never delete its academic record.  A school is a
-- business boundary with a lifecycle, not a disposable parent row: deleting it
-- would cascade memberships and orphan the evidence needed for a later reopen.
ALTER TABLE schools
    ADD COLUMN IF NOT EXISTS status VARCHAR(20) NOT NULL DEFAULT 'active'
        CHECK (status IN ('active', 'suspended', 'closed')),
    ADD COLUMN IF NOT EXISTS closed_at TIMESTAMP NULL,
    ADD COLUMN IF NOT EXISTS closed_by VARCHAR(255) NULL REFERENCES user_profiles(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS close_reason TEXT NULL;

CREATE INDEX IF NOT EXISTS idx_schools_status ON schools(status);
