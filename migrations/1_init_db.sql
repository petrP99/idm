-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS role
(
    id         bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name       text        NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS employee
(
    id         bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name       text        NOT NULL,
    role_id    bigint REFERENCES role (id),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE employee;
DROP TABLE role;
-- +goose StatementEnd