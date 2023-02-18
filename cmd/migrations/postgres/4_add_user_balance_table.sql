CREATE table if not exists "user_balance"
(
    "user_id"  int NOT null,
    "current"  int,
    "withdraw" int,
    CONSTRAINT fk_username FOREIGN KEY ("user_id") REFERENCES users ("user_id"),
    CONSTRAINT current_positive CHECK ( current >= 0 )
);
---- create above / drop below ----
DROP TABLE IF EXISTS "user_balance";
