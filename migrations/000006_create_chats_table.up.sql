CREATE TABLE IF NOT EXISTS chats
(
    id         BIGSERIAL PRIMARY KEY,
    buyer_id   BIGINT    NOT NULL,
    seller_id  BIGINT    NOT NULL,
    ad_id      BIGINT    NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (buyer_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (seller_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (ad_id) REFERENCES ads (id) ON DELETE CASCADE
);

CREATE Unique Index idx_chat_main On chats (buyer_id, seller_id, ad_id);
CREATE INDEX idx_buyer ON chats (buyer_id);
CREATE INDEX idx_seller ON chats (seller_id);
CREATE INDEX idx_ad ON chats (ad_id);
CREATE INDEX idx_created_at ON chats (created_at);