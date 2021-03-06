CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS "users" (
    "uid" uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    "email" VARCHAR NOT NULL UNIQUE,
    "password" VARCHAR NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);