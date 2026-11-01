BEGIN;

DROP INDEX IF EXISTS users.outbox_created_at_idx;
DROP INDEX IF EXISTS users.candidate_role_id_idx;
DROP INDEX IF EXISTS users.profile_role_id_idx;

DROP TABLE IF EXISTS users.outbox;
DROP TABLE IF EXISTS users.candidate;
DROP TABLE IF EXISTS users.profile;
DROP TABLE IF EXISTS users.role;

DROP SCHEMA IF EXISTS users;

COMMIT;