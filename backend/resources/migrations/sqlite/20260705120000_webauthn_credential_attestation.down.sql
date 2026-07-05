PRAGMA foreign_keys= OFF;
BEGIN;

CREATE TABLE webauthn_credentials_new AS
	SELECT id, created_at, updated_at, name, credential_id, public_key, attestation_type, transport, backup_eligible, backup_state, user_id
	FROM webauthn_credentials;

DROP TABLE webauthn_credentials;
ALTER TABLE webauthn_credentials_new RENAME TO webauthn_credentials;

COMMIT;
PRAGMA foreign_keys= ON;
