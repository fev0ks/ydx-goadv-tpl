CREATE table if not exists "orders"
(
    "number"      varchar PRIMARY KEY NOT null,
    "status"      varchar         NOT NULL,
    "accrual"     float4,
    "uploaded_at" varchar,
    "username"    varchar         NOT null,
    CONSTRAINT fk_username FOREIGN KEY ("username") REFERENCES users ("username")
);
---- create above / drop below ----
DROP TABLE IF EXISTS "orders";
