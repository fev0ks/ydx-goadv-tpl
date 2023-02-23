CREATE table if not exists "users"
(
    "user_id"  SERIAL  PRIMARY KEY NOT null,
    "username" varchar         NOT null,
    "password" bytea           NOT null
);
---- create above / drop below ----
DROP TABLE IF EXISTS "users";
