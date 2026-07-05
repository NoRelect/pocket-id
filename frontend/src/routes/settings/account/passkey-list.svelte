<script lang="ts">
	import { openConfirmDialog } from '$lib/components/confirm-dialog/';
	import GlassRowItem from '$lib/components/passkey-row.svelte';
	import * as Item from '$lib/components/ui/item/index.js';
	import { m } from '$lib/paraglide/messages';
	import WebauthnService from '$lib/services/webauthn-service';
	import type { MdsAuthenticator, Passkey } from '$lib/types/passkey.type';
	import { axiosErrorToast } from '$lib/utils/error-util';
	import { LucideKeyRound } from '@lucide/svelte';
	import { toast } from 'svelte-sonner';
	import RenamePasskeyModal from './rename-passkey-modal.svelte';

	let {
		passkeys = $bindable(),
		mdsMap = new Map<string, MdsAuthenticator>()
	}: {
		passkeys: Passkey[];
		mdsMap?: Map<string, MdsAuthenticator>;
	} = $props();

	const webauthnService = new WebauthnService();

	let passkeyToRename: Passkey | null = $state(null);

	async function deletePasskey(passkey: Passkey) {
		openConfirmDialog({
			title: m.delete_passkey_name({ passkeyName: passkey.name }),
			message: m.are_you_sure_you_want_to_delete_this_passkey(),
			confirm: {
				label: m.delete(),
				destructive: true,
				action: async () => {
					try {
						await webauthnService.removeCredential(passkey.id);
						passkeys = await webauthnService.listCredentials();
						toast.success(m.passkey_deleted_successfully());
					} catch (e) {
						axiosErrorToast(e);
					}
				}
			}
		});
	}
</script>

<Item.Group class="mt-3">
	{#each passkeys as passkey}
		<GlassRowItem
			label={passkey.name}
			description={m.added_on() + ' ' + new Date(passkey.createdAt).toLocaleDateString()}
			icon={LucideKeyRound}
			isCompromised={passkey.isCompromised}
			aaguid={passkey.aaguid}
			mdsAuthenticator={mdsMap.get(passkey.aaguid)}
			onRename={() => (passkeyToRename = passkey)}
			onDelete={() => deletePasskey(passkey)}
		/>
	{/each}
</Item.Group>

<RenamePasskeyModal
	bind:passkey={passkeyToRename}
	callback={async () => (passkeys = await webauthnService.listCredentials())}
/>
