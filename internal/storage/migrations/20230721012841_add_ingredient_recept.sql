-- +goose Up
-- +goose StatementBegin
CREATE TABLE "content"."ingredient_recept"(
    ingredient_id INTEGER REFERENCES "content"."ingredient"(id),
    quantity VARCHAR(100),
    recept_id INTEGER REFERENCES "content"."recept"(id),
    CONSTRAINT ingredients_recepts_pk PRIMARY KEY(ingredient_id,recept_id) );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "content"."ingredient_recept";
-- +goose StatementEnd
