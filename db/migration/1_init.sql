-- +migrate Up
CREATE TABLE "users" (
  "id" varchar(50) UNIQUE PRIMARY KEY NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz,
  "full_name" varchar(255) NOT NULL,
  "phone" varchar(15) NOT NULL,
  "country" varchar(50) NOT NULL,
  "email" varchar(255),
  "password" varchar(255),
  "last_sign" timestamptz,
  "role" varchar(255) NOT NULL,
  "aud" varchar(255) NOT NULL,
  CONSTRAINT unique_aud_phone UNIQUE ("aud", "phone"),
  CONSTRAINT unique_aud_email UNIQUE ("aud", "email")
);

CREATE TABLE "otp_authentications" (
  "id" SERIAL UNIQUE PRIMARY KEY,
  "aud" varchar(255) NOT NULL,
  "platform" varchar(255) NOT NULL,
  "phone" varchar(15) NOT NULL,
  "otp" varchar(10) NOT NULL,
  "entered_times" smallint NOT NULL DEFAULT 0,
  "resend_at" timestamptz NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz
);

CREATE TABLE oauth_tokens (
  "token_id" varchar(50) UNIQUE PRIMARY KEY NOT NULL,
  "aud" varchar(255) NOT NULL,
  "user_id" varchar(50) NOT NULL,
  "platform" varchar(255) NOT NULL,
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
DROP TABLE "oauth_tokens";
DROP TABLE "otp_authentications";
DROP TABLE "users";
