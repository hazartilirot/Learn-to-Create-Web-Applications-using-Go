BEGIN TRANSACTION;

drop table if exists users;

create table users (
    id serial primary key,
    age int,
    first_name varchar(150),
    last_name varchar(150),
    email varchar(150) unique not null,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
);

COMMIT;