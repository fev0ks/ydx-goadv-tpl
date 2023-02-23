CREATE table if not exists "user_withdraws"
(
    "user_id"     int NOT null,
    "withdraw_id" int NOT null,
    CONSTRAINT fk_username FOREIGN KEY ("user_id") REFERENCES users ("user_id"),
    CONSTRAINT fk_withdraw FOREIGN KEY ("withdraw_id") REFERENCES withdraws ("withdraw_id")
);
---- create above / drop below ----
DROP TABLE IF EXISTS "user_withdraws";
