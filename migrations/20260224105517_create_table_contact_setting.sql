-- +goose Up
-- +goose StatementBegin
CREATE TABLE im_contact.contact_setting (
    "id" uuid default uuidv7() primary key,
    "contact_id" uuid not null,
    "updated_at" timestamptz default now() not null,
    "allow_invites_from" int default 0 not null,
    CONSTRAINT "unique_contact_id" UNIQUE ("contact_id"),
    CONSTRAINT "fk_contact_id" FOREIGN KEY ("contact_id") REFERENCES im_contact.contact ("id") ON DELETE CASCADE
);


CREATE FUNCTION im_contact.create_setting_on_insert() RETURNS trigger AS $$
BEGIN
    INSERT INTO im_contact.contact_setting ("contact_id") VALUES (NEW.id);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER create_setting_on_contact_insert
AFTER INSERT ON im_contact.contact
FOR EACH ROW
EXECUTE FUNCTION im_contact.create_setting_on_insert();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS im_contact.contact_setting CASCADE;

DROP FUNCTION IF EXISTS im_contact.create_setting_on_insert() CASCADE;
-- +goose StatementEnd
