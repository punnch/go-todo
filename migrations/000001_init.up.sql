CREATE TABLE tasks (
    id             SERIAL                 PRIMARY KEY,
    title          VARCHAR(100)  NOT NULL CHECK(char_length(title) BETWEEN 1 AND 100),
    description    VARCHAR(1000) NOT NULL CHECK(char_length(description) BETWEEN 1 AND 100),
    completed      BOOLEAN       NOT NULL DEFAULT FALSE,
    created_at     TIMESTAMPTZ   NOT NULL
);