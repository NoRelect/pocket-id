<script lang="ts">
	import AuthenticatorDetailDialog from '$lib/components/authenticator-detail-dialog.svelte';
	import { Button } from '$lib/components/ui/button';
	import * as Item from '$lib/components/ui/item/index.js';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { m } from '$lib/paraglide/messages';
	import appConfigStore from '$lib/stores/application-configuration-store';
	import type { MdsAuthenticator } from '$lib/types/passkey.type';
	import { LucideAward, LucideCalendar, LucideInfo, LucidePencil, LucideTrash, LucideTriangleAlert, type Icon as IconType } from '@lucide/svelte';

	let {
		icon,
		onRename,
		onDelete,
		showRenameAction = true,
		label,
		description,
		isCompromised = false,
		attestationVerified = true,
		mdsAuthenticator
	}: {
		icon: typeof IconType;
		onRename?: () => void;
		onDelete: () => void;
		showRenameAction?: boolean;
		description?: string;
		label?: string;
		isCompromised?: boolean;
		attestationVerified?: boolean;
		mdsAuthenticator?: MdsAuthenticator;
	} = $props();

	let detailOpen = $state(false);
</script>

<Item.Root variant="transparent" class="hover:bg-muted transition-colors py-3 px-0 sm:px-4">
	<Item.Media class="bg-primary/10 text-primary rounded-full p-3">
		{#if icon}{@const Icon = icon}
			<Icon class="size-5" />
		{/if}
	</Item.Media>
	<Item.Content class="gap-0.5">
		<Item.Title class="flex items-center gap-2">
			{label}
			{#if isCompromised}
				<Tooltip.Provider>
					<Tooltip.Root>
						<Tooltip.Trigger>
							<span class="bg-destructive/10 text-destructive flex items-center gap-1 rounded px-1.5 py-0.5 text-xs font-medium">
								<LucideTriangleAlert class="size-3" />
								{m.passkey_compromised()}
							</span>
						</Tooltip.Trigger>
						<Tooltip.Content>{m.passkey_attestation_error_compromised()}</Tooltip.Content>
					</Tooltip.Root>
				</Tooltip.Provider>
			{/if}
		</Item.Title>
		{#if description}
			<Item.Description class="flex items-center">
				<LucideCalendar class="mr-1 size-3" />
				{description}
			</Item.Description>
		{/if}
		{#if mdsAuthenticator}
			<Item.Description>
				<button
					type="button"
					class="text-primary hover:text-primary/80 flex min-w-0 max-w-full items-center text-left transition-colors cursor-pointer"
					onclick={() => (detailOpen = true)}
					aria-label={m.passkey_authenticator_details()}
				>
					<LucideAward class="mr-1 size-3 shrink-0" />
					<span class="truncate">{mdsAuthenticator.description}</span>
				</button>
			</Item.Description>
		{/if}
		{#if !attestationVerified && $appConfigStore.passkeyAttestationMode === 'required'}
			<Item.Description class="text-destructive flex items-center">
				<LucideTriangleAlert class="mr-1 size-3 shrink-0" />
				{m.passkey_attestation_missing()}
			</Item.Description>
		{/if}
	</Item.Content>
	<Item.Actions>
		{#if showRenameAction && onRename}
			<Tooltip.Provider>
				<Tooltip.Root>
					<Tooltip.Trigger>
						<Button
							onclick={onRename}
							size="icon"
							variant="ghost"
							class="size-8"
							aria-label={m.rename()}
						>
							<LucidePencil class="size-4" />
						</Button>
					</Tooltip.Trigger>
					<Tooltip.Content>{m.rename()}</Tooltip.Content>
				</Tooltip.Root>
			</Tooltip.Provider>
		{/if}

		<Tooltip.Provider>
			<Tooltip.Root>
				<Tooltip.Trigger>
					<Button
						onclick={onDelete}
						size="icon"
						variant="ghost"
						class="hover:bg-destructive/10 hover:text-destructive size-8"
						aria-label={m.delete()}
					>
						<LucideTrash class="size-4" />
					</Button>
				</Tooltip.Trigger>
				<Tooltip.Content>{m.delete()}</Tooltip.Content>
			</Tooltip.Root>
		</Tooltip.Provider>
	</Item.Actions>
</Item.Root>

{#if mdsAuthenticator}
	<AuthenticatorDetailDialog bind:open={detailOpen} authenticator={mdsAuthenticator} />
{/if}
