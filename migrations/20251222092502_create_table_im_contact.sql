-- +goose Up
-- +goose StatementBegin
create table im_contact.contact (
    "id" uuid default uuidv7() primary key,
    "domain_id" bigint not null,
    "created_at" timestamptz default now() not null,
    "updated_at" timestamptz default now() not null,
    "created_by" bigint,
    "updated_by" bigint,

    "issuer_id" text not null,
    "application_id" text not null,

    "type" text not null,
    "name" text not null default '',
    "username" text not null,
    "metadata" jsonb,

    constraint contact_username_not_empty check (trim(username) <> ''),
    constraint contact_username_issuer_unique unique (issuer_id, username)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table im_contact.contact;
-- +goose StatementEnd
