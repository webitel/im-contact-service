-- +goose Up
-- +goose StatementBegin
ALTER TABLE im_contact.contact ADD COLUMN phone TEXT NOT NULL DEFAULT '';

-- WhatsApp and Viber BM contacts are keyed by their phone number (subject).
UPDATE im_contact.contact
SET phone = '+' || subject_id
WHERE issuer_id IN ('whatsapp', 'viber_bm')
  AND subject_id ~ '^[1-9][0-9]{6,14}$';

-- Portal contacts carry the OIDC phone_number claim in metadata.
UPDATE im_contact.contact c
SET phone = n.phone
FROM (
    -- Keep digits only, treat a leading "00" as the international prefix.
    SELECT id, '+' || regexp_replace(d, '^00', '') AS phone
    FROM (
        SELECT id,
               regexp_replace(metadata ->> '.phone_number', '[^0-9]', '', 'g') AS d
        FROM im_contact.contact
        WHERE phone = ''
          AND coalesce(metadata ->> '.phone_number', '') <> ''
    ) s
) n
WHERE c.id = n.id
  AND n.phone ~ '^\+[1-9][0-9]{6,14}$';

CREATE INDEX contact_domain_phone_idx ON im_contact.contact (domain_id, phone) WHERE phone <> '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS im_contact.contact_domain_phone_idx;
ALTER TABLE im_contact.contact DROP COLUMN phone;
-- +goose StatementEnd
