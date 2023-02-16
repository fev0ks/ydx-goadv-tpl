CREATE table if not exists "orders"
(
    "order_id"    int PRIMARY KEY NOT null,
    "status"      varchar         NOT NULL,
    "accrual"     int,
    "uploaded_at" timestamp
);
---- create above / drop below ----
DROP TABLE IF EXISTS "orders";
