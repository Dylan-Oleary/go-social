BEGIN;

ALTER TABLE users DROP COLUMN is_active;

DROP TABLE IF EXISTS user_invitations;

COMMIT;