CREATE TABLE IF NOT EXISTS appointments
(
    id           BIGSERIAL PRIMARY KEY,
    start_time   TIMESTAMP NOT NULL,
    duration     BIGINT    NOT NULL, -- в минутах
    location     TEXT      NOT NULL,
    ad_id        BIGINT    NOT NULL REFERENCES ads (id) ON DELETE CASCADE,
    buyer_id     BIGINT    NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    is_confirmed BOOLEAN   NOT NULL DEFAULT FALSE,
    is_canceled  BOOLEAN   NOT NULL DEFAULT FALSE,

    CONSTRAINT fk_buyer
        FOREIGN KEY (buyer_id)
            REFERENCES users (id)
            ON DELETE CASCADE,

    CONSTRAINT fk_ad
        FOREIGN KEY (ad_id)
            REFERENCES ads (id)
            ON DELETE CASCADE
);