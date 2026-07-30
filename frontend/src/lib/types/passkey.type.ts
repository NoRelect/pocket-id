export type Passkey = {
	id: string;
	name: string;
	createdAt: string;
	aaguid: string;
	attestationVerified: boolean;
	isCompromised: boolean;
};

export type MdsAuthenticator = {
	aaguid: string;
	description: string;
	isCompromised: boolean;
	fidoCertificationLevel: string;
	attestationTypes: string[];
	keyProtection: string[];
	attachmentHint: string[];
	icon: string;
	iconDark: string;
};
