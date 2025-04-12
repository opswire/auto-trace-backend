CREATE TABLE ads
(
    id              bigserial PRIMARY KEY,
    title           VARCHAR(255) NOT NULL,
    description     TEXT,
    price           FLOAT        NOT NULL,
    vin             VARCHAR(255) NOT NULL,
    is_token_minted BOOLEAN DEFAULT FALSE,
    brand           VARCHAR(255) NOT NULL,
    model           VARCHAR(255) NOT NULL,
    year_of_release SMALLINT     NOT NULL
);