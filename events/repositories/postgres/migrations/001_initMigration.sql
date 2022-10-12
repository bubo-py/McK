CREATE TABLE events (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    name VARCHAR(255) NOT NULL,
    startTime TIMESTAMP NOT NULL,
    endTime TIMESTAMP NOT NULL,
    description VARCHAR(510),
    alertTime TIMESTAMP
);

---- create above / drop below ----

drop table events;
