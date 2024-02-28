-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS goods (
    id SERIAL,
    project_id INT,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    priority SERIAL,
    removed BOOL NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id, project_id)
);

ALTER TABLE goods
ADD CONSTRAINT fk_project_constraint
FOREIGN KEY (project_id)
REFERENCES projects (id)
ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS goods;

ALTER TABLE goods
DROP CONSTRAINT fk_project_constraint;
-- +goose StatementEnd
