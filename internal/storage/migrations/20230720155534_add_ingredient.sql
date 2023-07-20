-- +goose Up
-- +goose StatementBegin

CREATE TABLE "content"."ingredient" (
    id  INTEGER PRIMARY KEY,
    name VARCHAR(100),
    description VARCHAR(255),
    href VARCHAR(255),
    parent_id INTEGER,
    updated_at timestamp with time zone NOT NULL
);
ALTER TABLE ONLY "content"."ingredient"
    ADD CONSTRAINT ingredient_fk_parent_id FOREIGN KEY (parent_id) 
    REFERENCES "content"."ingredient"(id) DEFERRABLE INITIALLY DEFERRED;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "content"."ingredient";
-- +goose StatementEnd