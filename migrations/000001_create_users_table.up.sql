CREATE TABLE users
(
    id       bigserial PRIMARY KEY,
    email    VARCHAR(100) NOT NULL,
    password VARCHAR(100) NOT NULL,
    role     VARCHAR(20)  NOT NULL DEFAULT 'user',
    name     VARCHAR(255) NOT NULL DEFAULT 'Аноним'
);

INSERT INTO users (email, password, role, name)
VALUES ('user1@example.com', '$2b$12$rid1ebIs4vD8lRzkEThokuOxZakuIIyu8RRx0mDtjClWLEf189bNa', 'user', 'Petya');
INSERT INTO users (email, password, role, name)
VALUES ('user2@example.com', '$2b$12$rid1ebIs4vD8lRzkEThokuOxZakuIIyu8RRx0mDtjClWLEf189bNa', 'user', 'Vanya');
INSERT INTO users (email, password, role, name)
VALUES ('user3@example.com', '$2b$12$rid1ebIs4vD8lRzkEThokuOxZakuIIyu8RRx0mDtjClWLEf189bNa', 'service', 'AutoTraceService');