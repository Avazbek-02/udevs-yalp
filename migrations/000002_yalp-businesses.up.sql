CREATE TABLE "businesses" (
  "id" UUID UNIQUE PRIMARY KEY NOT NULL,
  "owner_id" UUID,
  "name" VARCHAR(255) UNIQUE NOT NULL,
  "description" TEXT,
  "category" VARCHAR(255),
  "address" TEXT,
  "contact_info" TEXT,
  "photos" TEXT,
  "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);