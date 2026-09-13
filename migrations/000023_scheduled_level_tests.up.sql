-- NULL means the sheet has no scheduled date and can be taken at any time.
ALTER TABLE practice_sheets ADD COLUMN IF NOT EXISTS scheduled_at TIMESTAMP;

CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id VARCHAR(255) NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(200) NOT NULL,
    body TEXT,
    -- Points at whatever the notification is about, so the client can link to it.
    resource_type VARCHAR(50),
    resource_id UUID,
    -- When the referenced event happens; NULL for notifications with no date.
    scheduled_at TIMESTAMP,
    read_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notifications_user_created ON notifications(user_id, created_at DESC);
-- Rescheduling replaces the previous notification for the same student and sheet
-- instead of stacking a new one on every edit.
CREATE UNIQUE INDEX IF NOT EXISTS idx_notifications_user_resource
    ON notifications(user_id, type, resource_id)
    WHERE resource_id IS NOT NULL;
