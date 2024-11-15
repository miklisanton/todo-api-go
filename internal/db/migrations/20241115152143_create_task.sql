-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE task (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    due_date TIMESTAMP,
    overdue BOOLEAN DEFAULT FALSE,
    completed BOOLEAN DEFAULT FALSE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE task;
-- +goose StatementEn
