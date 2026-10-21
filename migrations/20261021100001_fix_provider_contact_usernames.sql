-- +goose Up
-- +goose StatementBegin

-- Facebook contacts were created with name = first name and
-- username = last name; the username must be the PSID (unique per issuer).
UPDATE im_contact.contact
SET name       = trim(name || ' ' || username),
    username   = subject_id,
    updated_at = now()
WHERE issuer_id = 'facebook'
  AND username <> subject_id;

-- WhatsApp contacts were created with username = profile name; the
-- username must be the wa_id (unique per issuer).
UPDATE im_contact.contact
SET username   = subject_id,
    updated_at = now()
WHERE issuer_id = 'whatsapp'
  AND username <> subject_id;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Data fix: the original usernames are not recoverable.
SELECT 1;
-- +goose StatementEnd
