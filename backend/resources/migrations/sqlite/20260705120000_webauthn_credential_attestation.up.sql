PRAGMA foreign_keys= OFF;
BEGIN;

ALTER TABLE webauthn_credentials
	ADD COLUMN aaguid TEXT NOT NULL DEFAULT '';

ALTER TABLE webauthn_credentials
	ADD COLUMN attestation_object BLOB;

COMMIT;
PRAGMA foreign_keys= ON;
