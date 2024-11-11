-- +goose Up
CREATE TABLE users (
    id uuid PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email TEXT NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL DEFAULT 'unset'
);

CREATE TABLE chirps (
    id uuid PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    body TEXT NOT NULL,
    user_id uuid NOT NULL ,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose Down
DROP TABLE chirps; 
DROP TABLE users;