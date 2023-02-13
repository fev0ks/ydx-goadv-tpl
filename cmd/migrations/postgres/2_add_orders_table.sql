CREATE table if not exists "orders"
(
    "order_id" int PRIMARY KEY NOT null,
    "username" varchar         NOT null,
    CONSTRAINT fk_username FOREIGN KEY ("username") REFERENCES users ("username")
);
---- create above / drop below ----
DROP TABLE IF EXISTS "orders";
