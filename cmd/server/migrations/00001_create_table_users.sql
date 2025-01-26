-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id bigserial primary key,
    name text,
    balance numeric(8,2) default 0,
    created_at timestamp with time zone default now()
);
insert into users(name) values ('Foo');
insert into users(name) values ('Bar');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
