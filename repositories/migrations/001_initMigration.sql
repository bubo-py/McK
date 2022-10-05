create table events (
    id serial primary key not null,
    name varchar(255) not null,
    startTime timestamp not null,
    endTime timestamp not null,
    description varchar(510),
    alertTime timestamp
);

---- create above / drop below ----

drop table events;

