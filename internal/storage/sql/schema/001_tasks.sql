-- +goose Up
CREATE TABLE tasks (
    id UUID PRIMARY KEY,
    description TEXT NOT NULL,
    is_completed BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    due_date TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE tasks;