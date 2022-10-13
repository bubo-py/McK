CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    login VARCHAR(30),
    password VARCHAR(60),
    timezone TEXT
);

---- create above / drop below ----

drop table events;
