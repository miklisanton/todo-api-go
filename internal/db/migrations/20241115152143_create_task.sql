-- +goose Up
-- +goose StatementBegin
CREATE TABLE task (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    due_date DATETIME,
    overdue BOOLEAN DEFAULT 0,
    completed BOOLEAN DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE task;
-- +goose StatementEnd
