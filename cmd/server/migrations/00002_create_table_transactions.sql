-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS transaction_types (
    id serial primary key,
    name text
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_transaction_types_name on transaction_types using hash (name);

insert into transaction_types(name) values('deposit');
insert into transaction_types(name) values('transfer');

CREATE TABLE IF NOT EXISTS transactions (
    id bigserial primary key,
    transaction_type_id integer REFERENCES transaction_types(id),
    user_id bigint not null REFERENCES users(id),
    amount numeric(8,2) not null,
    comment text,
    created_at timestamp with time zone default now()
);
CREATE INDEX IF NOT EXISTS idx_transactions_user_id on transactions(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transaction_types;
DROP TABLE IF EXISTS idx_transaction_types_name;
DROP TABLE IF EXISTS transactions;
DROP INDEX IF EXISTS idx_transactions_user_id;
-- +goose StatementEnd
