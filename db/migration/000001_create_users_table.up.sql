CREATE TABLE IF NOT EXISTS "users" (
  "id" text not null,
  "username" text not null unique,
  "email" text not null unique,
  "password" text not null,
  "bio" text,
  "image" text,
  "created_at" timestamptz not null default now(),
  "updated_at" timestamptz not null default now(),

  PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS "idx_users_username" ON "users" ("username");
CREATE INDEX IF NOT EXISTS "idx_users_email" ON "users" ("email");