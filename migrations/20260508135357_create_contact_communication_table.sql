-- +goose Up
create table if not exists "im_contact"."via"(
  "contact_id" uuid not null references "im_contact"."contact"("id") on delete cascade,
  "via" text not null check(trim("via")<>''),
  "disable" boolean not null default false,
  "disable_reason" text,
  "created_at" timestamp with time zone not null default now(),
  "updated_at" timestamp with time zone not null default now(),
  "metadata" jsonb,

  primary key ("contact_id", "via")
);

create or replace function "im_contact"."tg_via_integrity"()
returns trigger as $$
begin
  new.updated_at = now();
  new.created_at = old.created_at;
  return new;
end;
$$ language 'plpgsql';

create trigger "tg_via_integrity"
  before update on "im_contact"."via"
  for each row
  execute procedure "im_contact"."tg_via_integrity"();

-- +goose Down
drop trigger if exists "tg_via_integrity" on "im_contact"."communication";
drop function if exists "im_contact"."tg_via_integrity";
drop table if exists "im_contact"."communication";
