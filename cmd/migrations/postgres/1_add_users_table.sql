CREATE table if not exists "users"
(
    "username" varchar PRIMARY KEY NOT null,
    "password" bytea               NOT null
);
---- create above / drop below ----
DROP TABLE IF EXISTS "users";
