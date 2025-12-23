-- +goose Up
-- +goose StatementBegin
create schema im_contact;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop schema im_contact cascade;
-- +goose StatementEnd
