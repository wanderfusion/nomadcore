-- Write your migrate up statements here
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
