-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS logs (
    Id INT,
    ProjectId INT,
    Name VARCHAR(255) NOT NULL,
    Description TEXT,
    Priority INT,
    Removed BOOL NOT NULL,
    EventTime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS logs;
-- +goose StatementEnd
