-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA content;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA content;
-- +goose StatementEnd
