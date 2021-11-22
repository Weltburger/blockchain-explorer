CREATE DATABASE IF NOT EXISTS blockchain_explorer;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS "users" (
    "uid" uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    "email" VARCHAR NOT NULL UNIQUE,
    "password" VARCHAR NOT NULL,
    "changed_at" timestamptz NOT NULL DEFAULT('0001-01-01 00:00:00Z'),
    "created_at" timestamptz NOT NULL DEFAULT (now())
);