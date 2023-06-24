CREATE TABLE app_user_role
(id serial PRIMARY KEY,
name varchar(64) not null);

CREATE TABLE app_user
(id serial PRIMARY KEY,
app_user_role_id bigint not null REFERENCES app_user_role(id),
name varchar(64) not null,
password text not null,
last_update_date timestamp with time zone not null default CURRENT_TIMESTAMP
);

INSERT INTO app_user_role (id, name)
VALUES (1, 'system');

INSERT INTO app_user_role (id, name)
VALUES (2, 'user');

CREATE TABLE language
(id integer PRIMARY KEY,
name varchar(64) not null);

INSERT INTO language (id, name)
VALUES (1, 'English');

INSERT INTO language (id, name)
VALUES (2, 'Japanese');

CREATE TABLE genre
(id serial PRIMARY KEY,
name varchar(64) not null);

INSERT INTO genre (id, name)
VALUES (1, 'Programming');

CREATE TABLE author
(id serial PRIMARY KEY,
name varchar(256) not null);

CREATE TABLE book
(id serial PRIMARY KEY,
name varchar(64) not null,
genre_id bigint not null REFERENCES genre(id),
author_id bigint not null REFERENCES author(id),
language_id bigint not null REFERENCES language(id),
last_update_date timestamp with time zone not null default CURRENT_TIMESTAMP
);