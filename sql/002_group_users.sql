-- Write your migrate up statements here
create table group_users
(
    id         uuid      default gen_random_uuid() not null
        primary key,
    created_at timestamp default CURRENT_TIMESTAMP,
    updated_at timestamp default CURRENT_TIMESTAMP,
    group_id   uuid                                not null
        references groups,
    user_id    uuid                                not null
);

create unique index group_users_group_id_user_id_uindex
    on group_users (group_id, user_id);

---- create above / drop below ----

drop index group_users_group_id_user_id_uindex;
drop table group_users;
