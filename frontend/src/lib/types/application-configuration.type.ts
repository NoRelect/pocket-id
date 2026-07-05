import type { CustomClaim } from './custom-claim.type';

export type AppConfig = {
	appName: string;
	homePageUrl: string;
	allowOwnAccountEdit: boolean;
	allowUserSignups: 'disabled' | 'withToken' | 'open';
	emailOneTimeAccessAsUnauthenticatedEnabled: boolean;
	emailOneTimeAccessAsAdminEnabled: boolean;
	emailVerificationEnabled: boolean;
	ldapEnabled: boolean;
	disableAnimations: boolean;
	uiConfigDisabled: boolean;
	accentColor: string;
	requireUserEmail: boolean;
	tracingEnabled: boolean;
	passkeyAttestationMode: 'disabled' | 'optional' | 'required';
	passkeyAllowedAaguids: string;
	passkeyMinCertificationLevel: FidoCertificationLevel | '';
	passkeyRequiredKeyProtection: string;
};

export type FidoCertificationLevel =
	| 'FIDO_CERTIFIED'
	| 'FIDO_CERTIFIED_L1'
	| 'FIDO_CERTIFIED_L1plus'
	| 'FIDO_CERTIFIED_L2'
	| 'FIDO_CERTIFIED_L2plus'
	| 'FIDO_CERTIFIED_L3'
	| 'FIDO_CERTIFIED_L3plus';

export const fidoCertificationLevels: FidoCertificationLevel[] = [
	'FIDO_CERTIFIED',
	'FIDO_CERTIFIED_L1',
	'FIDO_CERTIFIED_L1plus',
	'FIDO_CERTIFIED_L2',
	'FIDO_CERTIFIED_L2plus',
	'FIDO_CERTIFIED_L3',
	'FIDO_CERTIFIED_L3plus'
];

export function formatCertificationLevel(level: string): string {
	return level
		.replace(/^FIDO_CERTIFIED$/, 'FIDO Certified')
		.replace(/^FIDO_CERTIFIED_L(\d)plus$/, 'FIDO Certified L$1+')
		.replace(/^FIDO_CERTIFIED_L(\d)$/, 'FIDO Certified L$1');
}

export type AllAppConfig = AppConfig & {
	// General
	sessionDuration: number;
	emailsVerified: boolean;
	signupDefaultUserGroupIDs: string[];
	signupDefaultCustomClaims: CustomClaim[];
	// Email
	smtpHost: string;
	smtpPort: string;
	smtpFrom: string;
	smtpUser: string;
	smtpPassword: string;
	smtpTls: 'none' | 'starttls' | 'tls';
	smtpSkipCertVerify: boolean;
	emailLoginNotificationEnabled: boolean;
	emailApiKeyExpirationEnabled: boolean;
	// LDAP
	ldapUrl: string;
	ldapBindDn: string;
	ldapBindPassword: string;
	ldapBase: string;
	ldapUserSearchFilter: string;
	ldapUserGroupSearchFilter: string;
	ldapSkipCertVerify: boolean;
	ldapAttributeUserUniqueIdentifier: string;
	ldapAttributeUserUsername: string;
	ldapAttributeUserEmail: string;
	ldapAttributeUserFirstName: string;
	ldapAttributeUserLastName: string;
	ldapAttributeUserDisplayName: string;
	ldapAttributeUserProfilePicture: string;
	ldapAttributeGroupMember: string;
	ldapAttributeGroupUniqueIdentifier: string;
	ldapAttributeGroupName: string;
	ldapAdminGroupName: string;
	ldapSoftDeleteUsers: boolean;
};

export type AppConfigRawResponse = {
	key: string;
	type: string;
	value: string;
}[];

export type AppVersionInformation = {
	isUpToDate: boolean | null;
	newestVersion: string | null;
	currentVersion: string;
};
