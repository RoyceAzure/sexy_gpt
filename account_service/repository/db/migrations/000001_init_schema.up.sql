CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- SQL dump generated using DBML (dbml-lang.org)
-- Database: PostgreSQL
-- Generated at: 2024-01-21T07:01:15.331Z

CREATE TYPE "gpt_key_type" AS ENUM (
  't3',
  't4'
);

CREATE TABLE "user" (
  "user_id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "user_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "is_email_verified" bool NOT NULL DEFAULT false,
  "hashed_password" varchar NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "sso_identifer" varchar,
  "is_internal" bool NOT NULL DEFAULT false,
  "cr_date" timestamptz NOT NULL DEFAULT (now()),
  "up_date" timestamptz,
  "cr_user" varchar NOT NULL,
  "up_user" varchar
);

CREATE TABLE "role" (
  "role_id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "role_name" varchar UNIQUE NOT NULL,
  "is_enable" bool NOT NULL DEFAULT true,
  "cr_date" timestamptz NOT NULL DEFAULT (now()),
  "up_date" timestamptz,
  "cr_user" varchar NOT NULL,
  "up_user" varchar
);

CREATE TABLE "user_role" (
  "user_id" uuid PRIMARY KEY NOT NULL,
  "role_id" uuid NOT NULL,
  "cr_date" timestamptz NOT NULL DEFAULT (now()),
  "up_date" timestamptz,
  "cr_user" varchar NOT NULL,
  "up_user" varchar
);

CREATE TABLE "vertify_email" (
  "id" bigserial PRIMARY KEY,
  "user_id" uuid NOT NULL,
  "email" varchar NOT NULL,
  "secret_code" varchar NOT NULL,
  "is_used" bool NOT NULL DEFAULT false,
  "cr_date" timestamptz NOT NULL DEFAULT (now()),
  "used_date" timestamptz,
  "expired_at" timestamptz NOT NULL DEFAULT (now() + interval '15 minutes')
);

CREATE TABLE "session" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "user_id" uuid UNIQUE NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" bool NOT NULL DEFAULT false,
  "cr_date" timestamptz NOT NULL DEFAULT (now()),
  "expired_at" timestamptz DEFAULT (now() + interval '3 days')
);

CREATE TABLE "gpt_key" (
  "key_id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "type" gpt_key_type NOT NULL,
  "expired_at" timestamptz NOT NULL,
  "assoicate_account" uuid,
  "max_usage" decimal NOT NULL,
  "current_usage" decimal NOT NULL,
  "max_share" decimal NOT NULL,
  "current_share" decimal NOT NULL
);

CREATE TABLE "account_key" (
  "user_id" uuid UNIQUE NOT NULL,
  "key_id" uuid NOT NULL,
  "expired_at" timestamptz NOT NULL,
  "cr_date" timestamptz NOT NULL DEFAULT (now()),
  "up_date" timestamptz
);

CREATE TABLE "msg_session" (
  "msg_session_id" uuid PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
  "user_id" uuid NOT NULL,
  "cr_date" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "msg" (
  "msg_id" bigserial NOT NULL,
  "msg_session_id" uuid NOT NULL,
  "user_msg" varchar NOT NULL,
  "response" varchar,
  "cr_date" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "user" ("user_id");

CREATE INDEX ON "role" ("role_name");

CREATE INDEX ON "user_role" ("user_id");

ALTER TABLE "user_role" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("user_id");

ALTER TABLE "vertify_email" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("user_id");

ALTER TABLE "session" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("user_id");

ALTER TABLE "gpt_key" ADD FOREIGN KEY ("assoicate_account") REFERENCES "user" ("user_id");

ALTER TABLE "account_key" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("user_id");

ALTER TABLE "account_key" ADD FOREIGN KEY ("key_id") REFERENCES "gpt_key" ("key_id");

ALTER TABLE "msg_session" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("user_id");

ALTER TABLE "msg" ADD FOREIGN KEY ("msg_session_id") REFERENCES "msg_session" ("msg_session_id");

