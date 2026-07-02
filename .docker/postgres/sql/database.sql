CREATE EXTENSION IF NOT EXISTS citext;


CREATE TABLE "roles" (
  "id" serial PRIMARY KEY,
  "name" text UNIQUE NOT NULL
);

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY DEFAULT uuidv7(),
  "email" citext NOT NULL,
  "name" text NOT NULL,
  "password" text NOT NULL,
  "role_id" int NOT NULL
);

CREATE TABLE "refresh_tokens" (
  "id" uuid PRIMARY KEY DEFAULT uuidv7(),
  "user_id" uuid NOT NULL,
  "token" text NOT NULL,
  "expires_at" timestamp NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL
);

CREATE TABLE "tipe_barang" (
  "id" serial PRIMARY KEY,
  "name" text NOT NULL
);

CREATE TABLE "satuan_barang" (
  "id" serial PRIMARY KEY,
  "satuan" text,
  "keterangan" text
);

CREATE TABLE "barang" (
  "kode" text PRIMARY KEY,
  "name" text NOT NULL,
  "tipe_id" integer NOT NULL,
  "satuan_id" integer NOT NULL,
  "quantity" integer DEFAULT 0,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL
);

CREATE UNIQUE INDEX "idx_users_email" ON "users" ("email");

CREATE INDEX "idx_users_name" ON "users" ("name");

CREATE INDEX "idx_barang_name" ON "barang" ("name");

CREATE INDEX "idx_tipe_barang" ON "barang" ("tipe_id");

ALTER TABLE "users"
ADD CONSTRAINT "fk_users_role"
FOREIGN KEY ("role_id") REFERENCES "roles" ("id")
ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE "barang"
ADD CONSTRAINT "fk_barang_tipe"
FOREIGN KEY ("tipe_id") REFERENCES "tipe_barang" ("id")
ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE "barang"
ADD CONSTRAINT "fk_barang_satuan"
FOREIGN KEY ("satuan_id") REFERENCES "satuan_barang" ("id")
ON DELETE RESTRICT ON UPDATE CASCADE;
