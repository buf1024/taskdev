create table task (
    id            integer primary key autoincrement,
    groupname     varchar(64),
    name          varchar(64),
    state         integer,
    outbin        varchar(256),
    log           varchar(256)
);