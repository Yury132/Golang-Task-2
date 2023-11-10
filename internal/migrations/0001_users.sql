-- +goose Up
create table if not exists public.uploads_info
(
    id        serial primary key,
    name      varchar(500),
    type      varchar(255),
    width     int,
    height    int,
    upload_at timestamp default now()
);

create table if not exists public.mini_info
(
    id        serial primary key,
    name      varchar(500),
    type      varchar(255),
    width     int,
    height    int,
    upload_at timestamp default now()
);

-- +goose Down
drop table public.uploads_info;

drop table public.mini_info;

