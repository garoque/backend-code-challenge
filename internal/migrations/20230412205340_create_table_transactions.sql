-- +goose Up
-- +goose StatementBegin
CREATE TABLE snapfi.transactions(
    id VARCHAR(36) NOT NULL UNIQUE,
    id_source VARCHAR(36) NOT NULL DEFAULT "",
    id_destination VARCHAR(36) NOT NULL,
    amount DECIMAL(9, 2) NOT NULL,
    state SMALLINT NOT NULL,
    created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    PRIMARY KEY (id)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE snapfi.transactions;
-- +goose StatementEnd
