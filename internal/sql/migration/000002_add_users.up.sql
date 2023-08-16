CREATE TABLE users
(
    "username"   varchar PRIMARY KEY,
    "name"       varchar     NOT NULL,
    "email"      varchar     NOT NULL UNIQUE,
    "password"   varchar     NOT NULL,
    "password_changed_at" timestamptz NOT NULL DEFAULT('0001-01-01 00:00:00Z'),
    "created_at" timestamptz NOT NULL DEFAULT (now())
);


ALTER TABLE "accounts"
    ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

ALTER TABLE "accounts"
    ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner", "currency");