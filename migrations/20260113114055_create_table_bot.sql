-- +goose Up
-- +goose StatementBegin
create table if not exists im_contact.bot (
    "id" uuid default uuidv7() primary key,
    "domain_id" bigint not null,
    "created_at" timestamptz default now() not null,
    "updated_at" timestamptz default now() not null,
    "flow_id" bigint not null,
    "display_name" text default '',
    constraint bot_domain_flow_unique unique(domain_id, flow_id)
);
create index if not exists idx_bot_domain_id on im_contact.bot using btree(domain_id, id);
create index if not exists idx_bot_flow_id on im_contact.bot using btree(flow_id);
--
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
drop table if exists im_contact.bot;
-- +goose StatementEnd