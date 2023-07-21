-- +goose Up
-- +goose StatementBegin
ALTER TABLE "content"."recept" ADD COLUMN "image" varchar;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "content"."recept" DROP COLUMN "image";
-- +goose StatementEnd
