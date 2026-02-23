CREATE TABLE IF NOT EXISTS user_profiles (
    user_id    VARCHAR(36) PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    numbers    JSONB NOT NULL DEFAULT '{}',
    strings    JSONB NOT NULL DEFAULT '{}',
    created_at BIGINT NOT NULL DEFAULT 0,
    updated_at BIGINT NOT NULL DEFAULT 0
);
