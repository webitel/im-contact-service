-- +goose Up
-- +goose StatementBegin
CREATE TABLE im_contact.contact_setting (
    "id" uuid default uuidv7() primary key,
    "contact_id" uuid not null,
    "rules" jsonb,
    "updated_at" timestamptz default now() not null
    CONSTRAINT "unique_contact_id" UNIQUE ("contact_id"),
    CONSTRAINT "fk_contact_id" FOREIGN KEY ("contact_id") REFERENCES im_contact.contact ("id") ON DELETE CASCADE
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE im_contact.contact_setting CASCADE;
-- +goose StatementEnd
