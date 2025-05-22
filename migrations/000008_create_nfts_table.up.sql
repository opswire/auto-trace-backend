CREATE TABLE IF NOT EXISTS nfts
(
    token_id     BIGSERIAL PRIMARY KEY,
    vin          VARCHAR(255) NOT NULL,
    metadata_url VARCHAR(511),
    is_minted    BOOLEAN               DEFAULT FALSE,
    created_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_vin ON nfts (vin);