-- +goose Up
-- +goose StatementBegin
CREATE TABLE snapfi.users(
    id VARCHAR(36) NOT NULL,
    name VARCHAR(80) NOT NULL,
    balance DECIMAL(9, 2) NOT NULL,
    created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    updated_at datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE snapfi.users;
-- +goose StatementEnd
