CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    login VARCHAR(30) UNIQUE NOT NULL,
    password VARCHAR(60) NOT NULL,
    timezone TEXT NOT NULL
);

---- create above / drop below ----

drop table events;
