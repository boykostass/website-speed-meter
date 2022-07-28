create table list
(
    id          uuid default gen_random_uuid() not null
        primary key,
    site        varchar(100)                   not null,
    date        varchar(100)                   not null,
    time        varchar(100)                   not null,
    delay       varchar(100)                   not null,
    performance varchar(100)                   not null
);

alter table list
    owner to postgres;

