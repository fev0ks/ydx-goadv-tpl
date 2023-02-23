CREATE table if not exists "user_orders"
(
    "user_id"  int  NOT null,
    "order_id" bigint  NOT null,
    CONSTRAINT fk_user_id FOREIGN KEY ("user_id") REFERENCES users ("user_id"),
    CONSTRAINT fk_order_id FOREIGN KEY ("order_id") REFERENCES orders ("order_id")
);
---- create above / drop below ----
DROP TABLE IF EXISTS "user_orders";
