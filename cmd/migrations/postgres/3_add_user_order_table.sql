CREATE table if not exists "user_orders"
(
    "username" varchar NOT null,
    "order"    varchar NOT null,
    CONSTRAINT fk_username FOREIGN KEY ("username") REFERENCES users ("username"),
    CONSTRAINT fk_order FOREIGN KEY ("order") REFERENCES orders ("number")
);
---- create above / drop below ----
DROP TABLE IF EXISTS "user_orders";
