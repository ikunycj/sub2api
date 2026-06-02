-- AI Studio: lightweight server-side conversations, messages, and temporary files.

CREATE TABLE IF NOT EXISTS ai_conversations (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    mode VARCHAR(16) NOT NULL CHECK (mode IN ('chat', 'image')),
    title TEXT NOT NULL DEFAULT '',
    pinned BOOLEAN NOT NULL DEFAULT FALSE,
    summary TEXT NOT NULL DEFAULT '',
    summary_message_id BIGINT,
    last_active_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ai_conversations_user_last_active
    ON ai_conversations(user_id, pinned DESC, last_active_at DESC, id DESC);

CREATE INDEX IF NOT EXISTS idx_ai_conversations_retention
    ON ai_conversations(last_active_at, id);

CREATE TABLE IF NOT EXISTS ai_messages (
    id BIGSERIAL PRIMARY KEY,
    conversation_id BIGINT NOT NULL REFERENCES ai_conversations(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(16) NOT NULL CHECK (role IN ('system', 'user', 'assistant')),
    kind VARCHAR(16) NOT NULL CHECK (kind IN ('chat', 'image')),
    content TEXT NOT NULL DEFAULT '',
    model TEXT NOT NULL DEFAULT '',
    status VARCHAR(16) NOT NULL DEFAULT 'completed',
    sequence BIGINT NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (conversation_id, sequence)
);

CREATE INDEX IF NOT EXISTS idx_ai_messages_conversation_sequence
    ON ai_messages(conversation_id, sequence);

CREATE INDEX IF NOT EXISTS idx_ai_messages_user_created
    ON ai_messages(user_id, created_at DESC, id DESC);

CREATE TABLE IF NOT EXISTS ai_attachments (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    conversation_id BIGINT REFERENCES ai_conversations(id) ON DELETE CASCADE,
    message_id BIGINT REFERENCES ai_messages(id) ON DELETE SET NULL,
    kind VARCHAR(32) NOT NULL CHECK (kind IN ('upload', 'generated_image')),
    filename TEXT NOT NULL DEFAULT '',
    storage_path TEXT NOT NULL,
    content_type TEXT NOT NULL DEFAULT '',
    byte_size BIGINT NOT NULL DEFAULT 0,
    sha256 TEXT NOT NULL DEFAULT '',
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ai_attachments_user_expires
    ON ai_attachments(user_id, expires_at, id);

CREATE INDEX IF NOT EXISTS idx_ai_attachments_expires
    ON ai_attachments(expires_at, id);

CREATE INDEX IF NOT EXISTS idx_ai_attachments_conversation
    ON ai_attachments(conversation_id, id);

CREATE INDEX IF NOT EXISTS idx_ai_attachments_message
    ON ai_attachments(message_id, id);
