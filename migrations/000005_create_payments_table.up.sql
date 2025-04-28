CREATE TABLE payments
(
    payment_id        BIGSERIAL PRIMARY KEY,
    user_id           BIGINT                   NOT NULL,
    ad_id             BIGINT                   NOT NULL,
    tariff_id         BIGINT                   NOT NULL,
    status            VARCHAR(20)              NOT NULL CHECK (
        status IN (
                   'pending',             -- Ожидает оплаты
                   'waiting_for_capture', -- Ожидает подтверждения
                   'succeeded',           -- Успешно завершен
                   'canceled'             -- Отмена оплаты
            )
        ),
    transaction_id    VARCHAR(255) UNIQUE,
    confirmation_link VARCHAR(255)             NOT NULL,
    expires_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at        TIMESTAMP                NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP                NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
            REFERENCES users (id)
            ON DELETE CASCADE,

    CONSTRAINT fk_ad
        FOREIGN KEY (ad_id)
            REFERENCES ads (id)
            ON DELETE CASCADE,

    CONSTRAINT fk_tariff
        FOREIGN KEY (tariff_id)
            REFERENCES tariffs (tariff_id)
            ON DELETE CASCADE
);

CREATE INDEX idx_payments_user ON payments (user_id);
CREATE INDEX idx_payments_ad ON payments (ad_id);
CREATE INDEX idx_payments_status ON payments (status);

CREATE OR REPLACE FUNCTION update_payments_modified_column()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_payments_updated_at
    BEFORE UPDATE
    ON payments
    FOR EACH ROW
EXECUTE FUNCTION update_payments_modified_column();