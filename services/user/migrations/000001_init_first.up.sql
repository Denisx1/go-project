BEGIN;
CREATE SCHEMA IF NOT EXISTS users;
CREATE TABLE users.role (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  name VARCHAR(255) NOT NULL UNIQUE,
  description TEXT NOT NULL,
  permissions TEXT NOT NULL,
  created_at TIMESTAMPTZ(3) NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ(3) NOT NULL DEFAULT now()
);
CREATE TABLE users.profile (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  version BIGINT NOT NULL DEFAULT 1,
  email VARCHAR(255) NOT NULL UNIQUE,
  user_name VARCHAR(255) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  is_email_verified BOOLEAN NOT NULL DEFAULT FALSE,
  is_active BOOLEAN NOT NULL DEFAULT FALSE,
  token_version INT NOT NULL DEFAULT 0,
  phone_number VARCHAR(50),
  role_id UUID NOT NULL,
  created_at TIMESTAMPTZ(3) NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ(3) NOT NULL DEFAULT now(),
  CONSTRAINT fk_user_role FOREIGN KEY (role_id) REFERENCES users.role (id) ON DELETE RESTRICT
);
CREATE TABLE users.candidate (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  version BIGINT NOT NULL DEFAULT 1,
  email VARCHAR(255) NOT NULL UNIQUE,
  user_name VARCHAR(255) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  action_token TEXT,
  role_id UUID NOT NULL,
  CONSTRAINT fk_candidate_role FOREIGN KEY (role_id) REFERENCES users.role (id) ON DELETE RESTRICT
);
CREATE TABLE users.outbox (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  aggregate_type TEXT NOT NULL,
  aggregate_id TEXT NOT NULL,
  event_type TEXT NOT NULL,
  payload JSONB NOT NULL,
  request_id TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX profile_role_id_idx ON users.profile (role_id);
CREATE INDEX candidate_role_id_idx ON users.candidate (role_id);
CREATE INDEX outbox_created_at_idx ON users.outbox (created_at);
COMMIT