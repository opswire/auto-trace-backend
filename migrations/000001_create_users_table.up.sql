CREATE TABLE users
(
    id       bigserial PRIMARY KEY,
    email    VARCHAR(100) NOT NULL,
    password VARCHAR(100) NOT NULL,
    role     VARCHAR(20)  NOT NULL DEFAULT 'user'
);