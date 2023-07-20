-- +goose Up
-- +goose StatementBegin
CREATE TABLE "content"."recept" (
    id  INTEGER PRIMARY KEY,
    name VARCHAR(100),
    description TEXT,
    cooking_time VARCHAR(100),
    number_servings VARCHAR(100),    
    image_src VARCHAR(255));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "content"."recept";
-- +goose StatementEnd