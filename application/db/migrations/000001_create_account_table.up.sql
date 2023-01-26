create table if not exists public.account
(
    id                 bigint                   not null
        primary key,
    phone              varchar(255),
    photo              varchar(255),
    about              varchar(255),
    city               varchar(255),
    country            varchar(255),
    status_code        varchar(255),
    reg_date           timestamp with time zone not null,
    message_permission varchar(255),
    last_online_time   timestamp with time zone,
    is_online          boolean,
    is_blocked         boolean                  not null,
    photo_id           varchar(255),
    photo_name         varchar(255),
    created            timestamp with time zone not null,
    updated            timestamp with time zone not null
);

alter table public.person
    owner to postgres;

