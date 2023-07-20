-- +goose Up
-- +goose StatementBegin

CREATE TABLE "content"."category" (
    slug VARCHAR(100) PRIMARY KEY,
    name VARCHAR(255),
    href VARCHAR(255),
    parent_slug VARCHAR(100)
);
ALTER TABLE ONLY "content"."category"
    ADD CONSTRAINT category_fk_parent_slug FOREIGN KEY (parent_slug) 
    REFERENCES "content"."category"(slug) DEFERRABLE INITIALLY DEFERRED;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "content"."category";
-- +goose StatementEnd
