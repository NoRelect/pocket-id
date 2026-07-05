ALTER TABLE webauthn_credentials
	ADD COLUMN aaguid TEXT NOT NULL DEFAULT '',
	ADD COLUMN attestation_object BYTEA;
