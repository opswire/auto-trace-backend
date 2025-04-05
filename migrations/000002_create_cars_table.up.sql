CREATE TABLE cars
(
    id              bigserial PRIMARY KEY,
    vin             VARCHAR(255) NOT NULL,
    is_token_minted BOOLEAN DEFAULT FALSE,
    brand           VARCHAR(255) NOT NULL,
    model           VARCHAR(255) NOT NULL,
    year_of_release SMALLINT     NOT NULL,
    image_url       VARCHAR(255)
);

CREATE UNIQUE INDEX ux_cars_vin ON cars (vin);