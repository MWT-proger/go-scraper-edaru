-- +goose Up
-- +goose StatementBegin
CREATE TABLE "content"."recept_category"(
    recept_id INTEGER REFERENCES "content"."recept"(id),
    category_slug VARCHAR(100) REFERENCES "content"."category"(slug),
    CONSTRAINT categories_recepts_pk PRIMARY KEY(category_slug,recept_id) );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "content"."recept_category";
-- +goose StatementEnd
