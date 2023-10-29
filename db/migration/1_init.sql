-- +migrate Up
CREATE TABLE "users" (
  "id" varchar(50) UNIQUE PRIMARY KEY NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz,
  "full_name" varchar(255) NOT NULL,
  "phone" varchar(15) UNIQUE NOT NULL,
  "country" varchar(50) NOT NULL,
  "email" varchar(100) UNIQUE,
  "birthday" timestamptz,
  "avatar" varchar(255),
  "password" varchar(255),
  "address" varchar(255),
  "gender" smallint,
  "status" smallint NOT NULL DEFAULT (1) -- 1: draft | 2: active
);

CREATE TABLE "otp_authentications" (
  "id" SERIAL UNIQUE PRIMARY KEY,
  "phone" varchar(15) NOT NULL,
  "otp" varchar(10) NOT NULL,
  "resend_at" timestamptz NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL,
  "deleted_at" timestamptz
);

CREATE TABLE oauth_tokens (
  "token_id" varchar(50) UNIQUE PRIMARY KEY NOT NULL,
  "user_id" varchar(50),
  "device_id" VARCHAR(255),
  "access_token" text NOT NULL,
  "refresh_token" text NOT NULL,
  "access_expires_at" timestamptz NOT NULL,
  "refresh_expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz,
  FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +migrate Down
DROP TABLE "users";
DROP TABLE "otp_authentications";
