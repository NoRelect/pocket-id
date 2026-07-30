<script lang="ts">
	import AuthenticatorDetailDialog from '$lib/components/authenticator-detail-dialog.svelte';
	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import { Button } from '$lib/components/ui/button';
	import * as Field from '$lib/components/ui/field';
	import { Input } from '$lib/components/ui/input';
	import * as Select from '$lib/components/ui/select';
	import { m } from '$lib/paraglide/messages';
	import WebAuthnService from '$lib/services/webauthn-service';
	import appConfigStore from '$lib/stores/application-configuration-store';
	import type { MdsAuthenticator } from '$lib/types/passkey.type';
	import {
		fidoCertificationLevels,
		formatCertificationLevel,
		type FidoCertificationLevel
	} from '$lib/types/application-configuration.type';
	import type { AllAppConfig } from '$lib/types/application-configuration.type';
	import { preventDefault } from '$lib/utils/event-util';
	import { onMount } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { LucideInfo } from '@lucide/svelte';

	let {
		callback,
		appConfig
	}: {
		appConfig: AllAppConfig;
		callback: (appConfig: Partial<AllAppConfig>) => Promise<void>;
	} = $props();

	const webAuthnService = new WebAuthnService();

	let attestationMode = $state<'disabled' | 'optional' | 'required'>(
		appConfig.passkeyAttestationMode === 'disabled' ||
			appConfig.passkeyAttestationMode === 'required'
			? appConfig.passkeyAttestationMode
			: 'optional'
	);

	let selectedAAGUIDs = $state<string[]>(
		appConfig.passkeyAllowedAaguids
			? appConfig.passkeyAllowedAaguids
					.split(',')
					.map((s) => s.trim())
					.filter(Boolean)
			: []
	);

	let minCertificationLevel = $state<FidoCertificationLevel | ''>(
		appConfig.passkeyMinCertificationLevel ?? ''
	);

	let mdsAuthenticators = $state<MdsAuthenticator[]>([]);
	let mdsLoading = $state(false);
	let isLoading = $state(false);
	let search = $state('');
	let showHidden = $state(false);
	let detailAuthenticator = $state<MdsAuthenticator | null>(null);
	let detailOpen = $state(false);

	let pruneOpen = $state(false);
	let pruneAffected = $state<MdsAuthenticator[]>([]);
	let pruneApply: (() => void) | null = null;
	let pruneRevert: (() => void) | null = null;

	onMount(async () => {
		mdsLoading = true;
		try {
			mdsAuthenticators = await webAuthnService.listMdsAuthenticators();
		} catch {
		} finally {
			mdsLoading = false;
		}
	});

	function openDetail(auth: MdsAuthenticator) {
		detailAuthenticator = auth;
		detailOpen = true;
	}

	function toggleAAGUID(aaguid: string) {
		if (selectedAAGUIDs.includes(aaguid)) {
			selectedAAGUIDs = selectedAAGUIDs.filter((a) => a !== aaguid);
		} else {
			selectedAAGUIDs = [...selectedAAGUIDs, aaguid];
		}
	}

	function clearSelection() {
		selectedAAGUIDs = [];
	}

	function supportsBasicFullAttestation(a: MdsAuthenticator): boolean {
		return a.attestationTypes.includes('basic_full');
	}

	function selectAll() {
		selectedAAGUIDs = mdsAuthenticators
			.filter((a) => !a.isCompromised && supportsBasicFullAttestation(a) && passesRestrictions(a))
			.map((a) => a.aaguid);
	}

	async function onSubmit() {
		isLoading = true;
		await callback({
			passkeyAttestationMode: attestationMode,
			passkeyAllowedAaguids: selectedAAGUIDs.join(','),
			passkeyMinCertificationLevel: minCertificationLevel
		}).finally(() => (isLoading = false));
		toast.success(m.passkey_attestation_updated_successfully());
	}

	const attestationModeOptions = [
		{
			value: 'disabled' as const,
			label: m.passkey_attestation_mode_disabled(),
			description: m.passkey_attestation_mode_disabled_description()
		},
		{
			value: 'optional' as const,
			label: m.passkey_attestation_mode_optional(),
			description: m.passkey_attestation_mode_optional_description()
		},
		{
			value: 'required' as const,
			label: m.passkey_attestation_mode_required(),
			description: m.passkey_attestation_mode_required_description()
		}
	];

	const currentModeLabel = $derived(
		attestationModeOptions.find((o) => o.value === attestationMode)?.label ?? attestationMode
	);

	const certificationLevelOptions = [
		{ value: '' as const, label: m.passkey_min_certification_level_none() },
		...fidoCertificationLevels.map((level) => ({
			value: level,
			label: formatCertificationLevel(level)
		}))
	];

	const currentCertificationLevelOption = $derived(
		certificationLevelOptions.find((o) => o.value === minCertificationLevel)
	);

	const showAllowList = $derived(attestationMode === 'required');

	$effect(() => {
		if (showAllowList) return;
		minCertificationLevel = '';
		selectedAAGUIDs = [];
	});

	function certificationLevelRank(level: string): number {
		const idx = fidoCertificationLevels.indexOf(level as FidoCertificationLevel);
		return idx === -1 ? 0 : idx + 1;
	}

	function passesRestrictions(a: MdsAuthenticator): boolean {
		if (minCertificationLevel) {
			if (
				certificationLevelRank(a.fidoCertificationLevel) <
				certificationLevelRank(minCertificationLevel)
			) {
				return false;
			}
		}
		return true;
	}

	let restrictionSnapshot = $state({ minCertificationLevel });

	$effect(() => {
		const changed = minCertificationLevel !== restrictionSnapshot.minCertificationLevel;

		if (!changed || pruneOpen) return;

		const affected = mdsAuthenticators.filter(
			(a) => selectedAAGUIDs.includes(a.aaguid) && !passesRestrictions(a)
		);

		if (affected.length === 0) {
			restrictionSnapshot = { minCertificationLevel };
			return;
		}

		const previousSnapshot = restrictionSnapshot;
		pruneAffected = affected;
		pruneApply = () => {
			const remove = new Set(affected.map((a) => a.aaguid));
			selectedAAGUIDs = selectedAAGUIDs.filter((aaguid) => !remove.has(aaguid));
			restrictionSnapshot = { minCertificationLevel };
		};
		pruneRevert = () => {
			minCertificationLevel = previousSnapshot.minCertificationLevel;
		};
		pruneOpen = true;
	});

	function resolvePrune(confirmed: boolean) {
		if (!pruneOpen) return;
		if (confirmed) {
			pruneApply?.();
		} else {
			pruneRevert?.();
		}
		pruneOpen = false;
		pruneAffected = [];
		pruneApply = null;
		pruneRevert = null;
	}

	const filteredAuthenticators = $derived(() => {
		const q = search.trim().toLowerCase();
		return mdsAuthenticators
			.filter(
				(a) =>
					!a.isCompromised &&
					supportsBasicFullAttestation(a) &&
					(!q || a.description.toLowerCase().includes(q) || a.aaguid.toLowerCase().includes(q))
			)
			.map((a) => ({ auth: a, excluded: !passesRestrictions(a) }))
			.filter((item) => showHidden || !item.excluded);
	});

	const hasHiddenAuthenticators = $derived(
		mdsAuthenticators.some(
			(a) => !a.isCompromised && supportsBasicFullAttestation(a) && !passesRestrictions(a)
		)
	);
</script>

<form onsubmit={preventDefault(onSubmit)}>
	<fieldset class="flex flex-col gap-5" disabled={$appConfigStore.uiConfigDisabled}>
		<Field.Field>
			<Field.Label>{m.passkey_attestation_mode()}</Field.Label>
			<Field.Description>{m.passkey_attestation_mode_description()}</Field.Description>
			<Select.Root
				type="single"
				value={attestationMode}
				onValueChange={(v) => (attestationMode = v as typeof attestationMode)}
			>
				<Select.Trigger class="w-full" aria-label={m.passkey_attestation_mode()}>
					{currentModeLabel}
				</Select.Trigger>
				<Select.Content>
					{#each attestationModeOptions as option}
						<Select.Item value={option.value}>
							<div class="flex flex-col items-start gap-1">
								<span class="font-medium">{option.label}</span>
								<span class="text-muted-foreground text-xs">
									{option.description}
								</span>
							</div>
						</Select.Item>
					{/each}
				</Select.Content>
			</Select.Root>
		</Field.Field>

		{#if showAllowList}
			<fieldset class="flex flex-col gap-5">
				<Field.Field>
					<Field.Label>{m.passkey_min_certification_level()}</Field.Label>
					<Field.Description>{m.passkey_min_certification_level_description()}</Field.Description>
					<Select.Root type="single" bind:value={minCertificationLevel}>
						<Select.Trigger class="w-full" aria-label={m.passkey_min_certification_level()}>
							{currentCertificationLevelOption?.label ?? m.passkey_min_certification_level_none()}
						</Select.Trigger>
						<Select.Content>
							{#each certificationLevelOptions as option}
								<Select.Item value={option.value}>
									{option.label}
								</Select.Item>
							{/each}
						</Select.Content>
					</Select.Root>
				</Field.Field>

				<Field.Field>
					<Field.Label>{m.passkey_allowed_aaguids()}</Field.Label>
					<Field.Description>{m.passkey_allowed_aaguids_description()}</Field.Description>

					<div class="mt-1 flex flex-wrap gap-2">
						<Button type="button" variant="outline" size="sm" onclick={clearSelection}>
							{m.passkey_attestation_preset_none()}
						</Button>
						<Button
							type="button"
							variant="outline"
							size="sm"
							onclick={selectAll}
							disabled={mdsLoading || mdsAuthenticators.length === 0}
						>
							{m.passkey_attestation_select_all()}
						</Button>
					</div>

					{#if !mdsLoading && mdsAuthenticators.length > 0}
						<Input
							type="search"
							placeholder="Search by name or AAGUID…"
							class="mt-2"
							bind:value={search}
						/>
					{/if}

					{#if !mdsLoading && hasHiddenAuthenticators}
						<label class="mt-2 flex cursor-pointer items-center gap-2 text-sm">
							<input type="checkbox" class="accent-primary shrink-0" bind:checked={showHidden} />
							{m.passkey_show_hidden_authenticators()}
						</label>
					{/if}

					<div class="mt-2 flex max-h-64 flex-col overflow-y-auto rounded-md border">
						{#if mdsLoading}
							<span class="text-muted-foreground p-3 text-sm">{m.loading()}</span>
						{:else if mdsAuthenticators.length === 0}
							<span class="text-muted-foreground p-3 text-sm">{m.passkey_mds_not_loaded()}</span>
						{:else}
							{@const visible = filteredAuthenticators()}
							{#if visible.length === 0}
								<span class="text-muted-foreground p-3 text-sm"
									>No authenticators match your search.</span
								>
							{:else}
								{#each visible as { auth, excluded }}
									<div
										class="flex items-center gap-3 border-b px-3 py-2 last:border-b-0 hover:bg-accent {excluded
											? 'opacity-50'
											: ''}"
									>
										<label
											class="flex min-w-0 flex-1 items-center gap-3 {excluded
												? 'cursor-not-allowed'
												: 'cursor-pointer'}"
										>
											<input
												type="checkbox"
												class="accent-primary shrink-0"
												checked={!excluded && selectedAAGUIDs.includes(auth.aaguid)}
												disabled={excluded}
												onchange={() => toggleAAGUID(auth.aaguid)}
											/>
											<span class="min-w-0 flex-1 text-sm">
												{auth.description || auth.aaguid}
											</span>
										</label>
										{#if auth.fidoCertificationLevel}
											<span
												class="bg-primary/10 text-primary shrink-0 rounded px-1.5 py-0.5 text-xs font-medium"
											>
												{formatCertificationLevel(auth.fidoCertificationLevel)}
											</span>
										{/if}
										<div
											role="button"
											tabindex="0"
											class="text-muted-foreground hover:text-foreground shrink-0 cursor-pointer"
											onclick={() => openDetail(auth)}
											onkeydown={(e) => (e.key === 'Enter' || e.key === ' ') && openDetail(auth)}
											aria-label={m.passkey_authenticator_details()}
										>
											<LucideInfo class="size-4" />
										</div>
									</div>
								{/each}
							{/if}
						{/if}
					</div>
				</Field.Field>
			</fieldset>
		{/if}

		<div>
			<Button type="submit" {isLoading}>{m.save()}</Button>
		</div>
	</fieldset>
</form>

<AuthenticatorDetailDialog bind:open={detailOpen} authenticator={detailAuthenticator} />

<AlertDialog.Root
	bind:open={pruneOpen}
	onOpenChange={(v) => {
		if (!v) resolvePrune(false);
	}}
>
	<AlertDialog.Content class="z-9999">
		<AlertDialog.Header>
			<AlertDialog.Title>{m.passkey_prune_selection_title()}</AlertDialog.Title>
			<AlertDialog.Description>
				{m.passkey_prune_selection_description()}
			</AlertDialog.Description>
		</AlertDialog.Header>

		<ul class="max-h-48 overflow-y-auto rounded-md border">
			{#each pruneAffected as auth}
				<li class="border-b px-3 py-2 text-sm last:border-b-0">
					{auth.description || auth.aaguid}
				</li>
			{/each}
		</ul>

		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={() => resolvePrune(false)}>{m.cancel()}</AlertDialog.Cancel>
			<AlertDialog.Action>
				{#snippet child()}
					<Button variant="destructive" onclick={() => resolvePrune(true)}>
						{m.passkey_prune_selection_confirm()}
					</Button>
				{/snippet}
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
