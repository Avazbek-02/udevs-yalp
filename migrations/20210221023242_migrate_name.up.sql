
CREATE TYPE "user_type" AS ENUM ('user', 'admin');
CREATE TYPE "user_role" AS ENUM ('user', 'admin', 'superadmin');
CREATE TYPE "gender" AS ENUM ('male', 'female');
CREATE TYPE "user_status" AS ENUM ('active', 'blocked', 'inverify');
CREATE TYPE "platform" AS ENUM ('admin', 'web', 'mobile');

CREATE TABLE IF NOT EXISTS "users" (
  "id" UUID UNIQUE PRIMARY KEY NOT NULL,
  "user_type" user_type NOT NULL,
  "user_role" user_role NOT NULL,
  "email" VARCHAR(255) UNIQUE NOT NULL,
  "password_hash" VARCHAR(255) NOT NULL,
  "username" VARCHAR(255) not null,
  "full_name" VARCHAR(255) not null,
  "gender" gender NOT NULL DEFAULT 'male',
  "profile_picture" VARCHAR(255),
  "bio" TEXT,
  "status" user_status NOT NULL DEFAULT 'inverify',
  "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE IF NOT EXISTS "sessions" (
  "id" UUID UNIQUE PRIMARY KEY NOT NULL,
  "user_id" UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  "user_agent" TEXT NOT NULL,
  "platform" platform NOT NULL,
  "ip_address" VARCHAR(64) NOT NULL,
  "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);
