-- +goose Up
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(50) UNIQUE NOT NULL,
    description VARCHAR(200) NOT NULL,
    completed BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL
);

INSERT INTO tasks (title, description, completed, created_at)
VALUES ('coding', '2 hours', FALSE, '2026-02-07 16:42:00');

-- +goose Down
DROP TABLE tasks;