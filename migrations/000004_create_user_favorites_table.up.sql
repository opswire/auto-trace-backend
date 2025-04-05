CREATE TABLE user_favorites
(
    user_id    INTEGER REFERENCES users (id) ON DELETE CASCADE,
    ad_id      INTEGER REFERENCES ads (id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, ad_id)
);
