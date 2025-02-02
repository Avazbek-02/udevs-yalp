
CREATE TYPE "user_type" AS ENUM ('user', 'admin','businessman');
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
  "avatar_id" VARCHAR(255) UNIQUE,
  "bio" TEXT DEFAULT ' ',
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
  "is_active" bool NOT NULL,
  "expires_at" timestamp,
  "last_active_at" timestamp,
  "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);

INSERT INTO "users" (
  "id", "user_type", "user_role", "email", "password_hash", 
  "username", "full_name", "gender","avatar_id", "status"
) VALUES (
  'e1ebed26-59e6-4eb2-bd35-13d504e79cd3',
  'admin',
  'superadmin',
  'apalonavalon@gmail.com',
  '$2a$10$Bmk0SXUCNjc/3gA3R4/srOvzHabpNx/WvgHgfPtVKSEu3.9TTT54e',
  'superadmin',
  'Default Super Admin',
  'male',
  'admin',
  'active'
);