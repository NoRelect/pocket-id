import UserService from '$lib/services/user-service';
import WebAuthnService from '$lib/services/webauthn-service';
import type { MdsAuthenticator } from '$lib/types/passkey.type';
import type { PageLoad } from './$types';

export const load: PageLoad = async () => {
	const webauthnService = new WebAuthnService();
	const userService = new UserService();

	const [account, passkeys, mdsAuthenticators] = await Promise.all([
		userService.getCurrent(),
		webauthnService.listCredentials(),
		webauthnService.listMdsAuthenticators().catch(() => [] as MdsAuthenticator[])
	]);

	const mdsMap = new Map(mdsAuthenticators.map((a) => [a.aaguid, a]));

	return {
		account,
		passkeys,
		mdsMap
	};
};
