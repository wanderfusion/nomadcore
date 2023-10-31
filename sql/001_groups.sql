-- Write your migrate up statements here
-- auto-generated definition
create table groups
(
    id          uuid      default gen_random_uuid() not null
        primary key,
    created_at  timestamp default CURRENT_TIMESTAMP,
    updated_at  timestamp default CURRENT_TIMESTAMP,
    user_id     uuid                                not null,
    name        text                                not null,
    description text
);

---- create above / drop below ----
drop table groups;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
