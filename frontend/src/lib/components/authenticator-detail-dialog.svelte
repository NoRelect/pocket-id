<script lang="ts">
	import CopyToClipboard from '$lib/components/copy-to-clipboard.svelte';
	import * as Dialog from '$lib/components/ui/dialog';
	import { m } from '$lib/paraglide/messages';
	import { formatCertificationLevel } from '$lib/types/application-configuration.type';
	import type { MdsAuthenticator } from '$lib/types/passkey.type';
	import { mode } from 'mode-watcher';

	let {
		open = $bindable(false),
		authenticator
	}: {
		open?: boolean;
		authenticator: MdsAuthenticator | null;
	} = $props();

	let contentRef = $state<HTMLElement | null>(null);

	const icon = $derived(
		authenticator
			? (mode.current === 'dark' && authenticator.iconDark) || authenticator.icon || null
			: null
	);
</script>

<Dialog.Root bind:open>
	<Dialog.Content
		class="max-w-xl"
		onOpenAutoFocus={(e) => {
			e.preventDefault();
			contentRef?.focus();
		}}
		bind:ref={contentRef}
		tabindex={-1}
	>
		{#if authenticator}
			<Dialog.Header class="min-w-0">
				<div class="flex min-w-0 items-center gap-4">
					{#if icon}
						<img
							src={icon}
							alt={authenticator.description}
							class="size-12 shrink-0 rounded object-contain"
						/>
					{/if}
					<div class="min-w-0">
						<Dialog.Title class="truncate"
							>{authenticator.description || authenticator.aaguid}</Dialog.Title
						>
						{#if authenticator.isCompromised}
							<span
								class="bg-destructive/10 text-destructive mt-1 inline-flex items-center rounded px-1.5 py-0.5 text-xs font-medium"
							>
								{m.passkey_authenticator_status_compromised()}
							</span>
						{/if}
					</div>
				</div>
			</Dialog.Header>

			<div class="mt-2 flex flex-col gap-3 text-sm">
				{@render Row(m.passkey_authenticator_aaguid(), Aaguid)}
				{#if authenticator.fidoCertificationLevel}
					{@render Row(m.passkey_authenticator_certification_level(), CertLevel)}
				{/if}
				{#if authenticator.keyProtection.length > 0}
					{@render Row(m.passkey_authenticator_key_protection(), KeyProtection)}
				{/if}
				{#if authenticator.attachmentHint.length > 0}
					{@render Row(m.passkey_authenticator_attachment_hint(), AttachmentHint)}
				{/if}
				{#if authenticator.attestationTypes.length > 0}
					{@render Row(m.passkey_authenticator_attestation_types(), AttestationTypes)}
				{/if}
			</div>
		{/if}
	</Dialog.Content>
</Dialog.Root>

{#snippet Row(label: string, children: import('svelte').Snippet)}
	<div class="flex items-start gap-3">
		<span class="text-muted-foreground w-36 shrink-0 pt-0.5">{label}</span>
		<span class="min-w-0 flex-1">{@render children()}</span>
	</div>
{/snippet}

{#snippet Aaguid()}
	{#if authenticator}
		<div class="min-w-0 overflow-x-auto">
			<CopyToClipboard value={authenticator.aaguid}>
				<span class="font-mono text-xs whitespace-nowrap">{authenticator.aaguid}</span>
			</CopyToClipboard>
		</div>
	{/if}
{/snippet}

{#snippet CertLevel()}
	<span class="bg-primary/10 text-primary rounded px-1.5 py-0.5 text-xs font-medium wrap-break-word">
		{authenticator ? formatCertificationLevel(authenticator.fidoCertificationLevel) : ''}
	</span>
{/snippet}

{#snippet KeyProtection()}
	{@render TagList(authenticator?.keyProtection ?? [])}
{/snippet}

{#snippet AttachmentHint()}
	{@render TagList(authenticator?.attachmentHint ?? [])}
{/snippet}

{#snippet AttestationTypes()}
	{@render TagList(authenticator?.attestationTypes ?? [])}
{/snippet}

{#snippet TagList(values: string[])}
	<div class="flex flex-wrap gap-1">
		{#each values as value}
			<span
				class="bg-muted text-muted-foreground rounded px-1.5 py-0.5 font-mono text-xs break-all"
			>
				{value}
			</span>
		{/each}
	</div>
{/snippet}
