CREATE table if not exists "withdraws"
(
    "withdraw_id"  SERIAL PRIMARY KEY NOT null,
    "order_id"     bigint                NOT null,
    "sum"          int                NOT null,
    "processed_at" timestamp
);
---- create above / drop below ----
DROP TABLE IF EXISTS "withdraws";
