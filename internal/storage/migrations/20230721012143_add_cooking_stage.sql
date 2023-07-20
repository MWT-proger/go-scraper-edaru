-- +goose Up
-- +goose StatementBegin
CREATE TABLE "content"."cooking_stage" (
    id  SERIAL PRIMARY KEY,
    number VARCHAR(100),
    description TEXT,
    recept_id INTEGER);
ALTER TABLE ONLY "content"."cooking_stage"
    ADD CONSTRAINT cooking_stage_fk_recept_id FOREIGN KEY (recept_id) 
    REFERENCES "content"."recept"(id) DEFERRABLE INITIALLY DEFERRED;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "content"."cooking_stage";
-- +goose StatementEnd