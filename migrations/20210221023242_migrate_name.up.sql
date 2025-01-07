CREATE TYPE "user_type" AS ENUM (
  'user',
  'admin'
);

CREATE TYPE "user_role" AS ENUM (
  'user',
  'admin',
  'superadmin'
);

CREATE TYPE "gender" AS ENUM (
  'male',
  'female'
);

CREATE TYPE "user_status" AS ENUM (
  'active',
  'blocked',
  'inverify'
);

CREATE TYPE "platform" AS ENUM (
  'admin',
  'web',
  'mobile'
);

CREATE TABLE if not exists "users" (
  "id" UUID UNIQUE PRIMARY KEY NOT NULL,
  "user_type" user_type NOT NULL,
  "user_role" user_role NOT NULL,
  "email" VARCHAR(255) UNIQUE NOT NULL,
  "password_hash" VARCHAR(255) NOT NULL,
  "name" VARCHAR(255),
  "gender" gender NOT NULL DEFAULT 'male',
  "profile_picture" VARCHAR(255),
  "bio" TEXT,
  "status" user_status NOT NULL DEFAULT 'inverify',
  "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE if not exists "sessions" (
  "id" UUID UNIQUE PRIMARY KEY NOT NULL,
  "user_id" UUID NOT NULL REFERENCES Users(id) ON delete cascade,
  "user_agent" TEXT NOT NULL,
  "platform" ENUM(admin,web,mobile) NOT NULL,
  "ip_address" VARCHAR(64) NOT NULL,
  "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);
