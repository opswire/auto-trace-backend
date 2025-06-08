CREATE TABLE IF NOT EXISTS messages
(
    id         BIGSERIAL PRIMARY KEY,
    chat_id    BIGINT    NOT NULL,
    sender_id  BIGINT    NOT NULL,
    text       TEXT      NOT NULL,
    is_read    BOOLEAN   NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    image_url VARCHAR(511) DEFAULT NULL,

    FOREIGN KEY (sender_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (chat_id) REFERENCES chats (id) ON DELETE CASCADE
);

CREATE INDEX idx_sender ON messages (sender_id);
CREATE INDEX idx_chat ON messages (chat_id);
