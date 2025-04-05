CREATE TABLE ads (
    id bigserial PRIMARY KEY,
    car_id INTEGER REFERENCES cars (id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    price FLOAT NOT NULL
);