-- Write your migrate up statements here
create table user_profiles
(
    id         uuid      default gen_random_uuid() not null
        primary key,
    created_at timestamp default CURRENT_TIMESTAMP,
    updated_at timestamp default CURRENT_TIMESTAMP,
    user_id    uuid                                not null
        unique,
    bio        text,
    interests  text,
    metadata   text
);

---- create above / drop below ----

drop table user_profiles;
