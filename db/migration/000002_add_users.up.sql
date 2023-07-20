CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE FUNCTION encrypt_data(data text) RETURNS bytea AS $$
DECLARE
  encrypted_data bytea;
BEGIN
  encrypted_data = pgcrypto.encrypt(data, 'safedata');
  RETURN encrypted_data;
END;
$$ LANGUAGE plpgsql;

CREATE FUNCTION decrypt_data(encrypted_data bytea) RETURNS text AS $$
DECLARE
  data text;
BEGIN
  data = pgcrypto.decrypt(encrypted_data, 'safedata');
  RETURN data;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "hashed_pass" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "email_encrypt" bytea,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

-- CREATE UNIQUE INDEX ON "accounts" ("owner", "currency");

ALTER TABLE "accounts" ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner", "currency");
