-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS hash
(
    id   bigserial PRIMARY KEY,
    hash varchar NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS hash;
-- +goose StatementEnd
