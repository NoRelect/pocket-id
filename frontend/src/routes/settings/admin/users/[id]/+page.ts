import UserService from '$lib/services/user-service';
import WebAuthnService from '$lib/services/webauthn-service';
import type { MdsAuthenticator } from '$lib/types/passkey.type';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params }) => {
	const userService = new UserService();
	const webauthnService = new WebAuthnService();
	const [user, passkeys, mdsAuthenticators] = await Promise.all([
		userService.get(params.id),
		userService.listUserPasskeys(params.id),
		webauthnService.listMdsAuthenticators().catch(() => [] as MdsAuthenticator[])
	]);

	const mdsMap = new Map(mdsAuthenticators.map((a) => [a.aaguid, a]));

	return {
		user,
		passkeys,
		mdsMap
	};
};
