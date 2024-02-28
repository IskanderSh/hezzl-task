-- +goose Up
-- +goose StatementBegin
INSERT INTO projects (name)
VALUES ('Первая запись');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM projects WHERE name='Начальное название';
-- +goose StatementEnd
