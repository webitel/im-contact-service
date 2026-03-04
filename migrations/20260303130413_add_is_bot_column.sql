-- +goose Up
-- +goose StatementBegin
ALTER TABLE im_contact.contact ADD COLUMN is_bot BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE im_contact.contact DROP COLUMN is_bot;
-- +goose StatementEnd
