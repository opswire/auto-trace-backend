CREATE TABLE tariffs
(
    tariff_id    SERIAL PRIMARY KEY,
    name         VARCHAR(255)   NOT NULL UNIQUE,
    description  TEXT,
    price        DECIMAL(12, 2) NOT NULL CHECK (price >= 0),
    currency     VARCHAR(3)     NOT NULL DEFAULT 'RUB',
    duration_min INT            NOT NULL CHECK (duration_min > 0), -- в минутах
    is_active    BOOLEAN        NOT NULL DEFAULT TRUE,
    created_at   TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tariffs_active ON tariffs (is_active);

-- Триггер для обновления updated_at
CREATE OR REPLACE FUNCTION update_tariffs_modified_column()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_tariffs_updated_at
    BEFORE UPDATE
    ON tariffs
    FOR EACH ROW
EXECUTE FUNCTION update_tariffs_modified_column();

-- Тестовые данные
INSERT INTO tariffs (name, description, price, duration_min)
VALUES ('Стандарт', 'Продвижение в общем списке на 24 часа', 490.00, 24 * 60),
       ('Премиум', 'Закрепление в топе на 7 дней + выделение цветом', 1490.00, 7 * 24 * 60);
