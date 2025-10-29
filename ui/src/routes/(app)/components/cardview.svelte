<script lang="ts">
	import CopyButton from '$lib/components/CopyButton.svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import {
		Power,
		PowerOff,
		RotateCcw,
		CirclePause,
		FileCog,
		Download,
		Trash,
		Undo2,
		Shredder
	} from '@lucide/svelte/icons';

	let {
		instances
	}: {
		instances: App.Instance[];
	} = $props();

	function handleAction(instId: string, action: string) {
		console.log(`Perform ${action} on`, instId);
	}

	function getHostname(hostname: string, port: number): string {
		let hn = hostname || '127.0.0.1';
		let pt = port || 6379;
		return hn + ':' + pt;
	}
</script>

<div class="grid gap-4 sm:grid-cols-1 xl:grid-cols-2 2xl:grid-cols-3">
	{#each instances as inst (inst.id)}
		<Card.Root class="min-w-sm border-muted/30 bg-muted/30 dark:border-muted/20 dark:bg-muted/10">
			<Card.Header>
				<div class="flex items-center justify-between">
					<div class="flex items-center gap-2">
						<Card.Title class="text-base font-medium tracking-tight">
							{inst.name}
						</Card.Title>

						<span
							class="
								rounded-full px-2 py-0.5 text-xs font-semibold capitalize
								transition-colors
								{inst.status === 'running'
								? 'bg-emerald-100 text-emerald-700 dark:bg-emerald-500/10 dark:text-emerald-400'
								: inst.status === 'stopped'
									? 'bg-yellow-100 text-yellow-700 dark:bg-yellow-500/10 dark:text-yellow-400'
									: inst.status === 'deleted'
										? 'bg-red-100 text-red-700 dark:bg-red-500/10 dark:text-red-400'
										: 'bg-gray-100 text-gray-700 dark:bg-gray-500/10 dark:text-gray-400'}
							"
						>
							{inst.status}
						</span>
					</div>
				</div>

				<Card.Description class="mt-1  truncate text-xs text-muted-foreground">
					{inst.id}
				</Card.Description>

				<Card.Action class="mt-3 flex flex-row gap-2">
					{#if inst.status === 'running'}
						<Button
							variant="secondary"
							size="icon"
							class="size-8 hover:text-red-500 dark:hover:text-red-400"
							onclick={() => handleAction(inst.id, 'stop')}
							title="Stop Instance"
						>
							<PowerOff />
						</Button>
						<Button
							variant="secondary"
							size="icon"
							class="size-8 hover:text-yellow-500 dark:hover:text-yellow-400"
							onclick={() => handleAction(inst.id, 'restart')}
							title="Restart Instance"
						>
							<RotateCcw />
						</Button>
					{:else if inst.status === 'stopped'}
						<Button
							variant="secondary"
							size="icon"
							class="size-8 hover:text-emerald-600 dark:hover:text-emerald-400"
							onclick={() => handleAction(inst.id, 'start')}
							title="Start Instance"
						>
							<Power />
						</Button>
					{:else if inst.status === 'deleted'}
						<Button
							variant="secondary"
							size="default"
							class=" hover:text-emerald-500 dark:hover:text-emerald-400"
							title="Restore Instance"
						>
							<Undo2 />Restore
						</Button>
					{:else}
						<Button
							variant="secondary"
							size="icon"
							class="size-8 text-muted-foreground hover:text-foreground/80"
							title="No actions available"
						>
							<CirclePause />
						</Button>
					{/if}
				</Card.Action>
			</Card.Header>

			<Card.Content class="flex flex-col py-2 text-sm text-muted-foreground">
				{#if inst.status !== 'deleted'}
					<div class="flex flex-row items-center gap-4">
						<span class="text-lg">{getHostname(inst.primaryHostname, inst.port)}</span>
						<CopyButton value={getHostname(inst.primaryHostname, inst.port)} size="sm" />
					</div>
				{:else}
					<span class="text-lg"
						>Can be restored until {new Date(
							new Date(inst.createdAt).getTime() + 30 * 24 * 60 * 60 * 1000
						).toLocaleString()}</span
					>
				{/if}
			</Card.Content>

			{#if inst.status !== 'deleted'}
				<Card.Footer class="flex-row justify-between text-xs text-muted-foreground">
					<span>{new Date(inst.createdAt).toLocaleString()}</span>
					<div class="flex flex-row gap-2">
						<Button
							variant="secondary"
							size="default"
							class="hover:text-neutral-600 dark:hover:text-neutral-300"
							onclick={() => handleAction(inst.id, 'start')}
							title="Open Details"
						>
							<FileCog />Details
						</Button>
						<Button
							variant="secondary"
							size="icon"
							class=" hover:text-red-500 dark:hover:text-red-400"
							title="Remove Instance"
						>
							<Trash />
						</Button>
					</div>
				</Card.Footer>
			{:else}
				<Card.Footer
					class="bottom-0 mt-auto flex-row justify-between text-xs text-muted-foreground"
				>
					<span>{new Date(inst.createdAt).toLocaleString()}</span>
					<div class="flex flex-row gap-2">
						<Button
							variant="secondary"
							size="default"
							class="hover:text-neutral-600 dark:hover:text-neutral-300"
							onclick={() => handleAction(inst.id, 'start')}
							title="Open Details"
						>
							<Download />Export
						</Button>
						<Button
							variant="secondary"
							size="icon"
							class=" hover:text-red-500 dark:hover:text-red-400"
							title="Remove Instance"
						>
							<Shredder />
						</Button>
					</div>
				</Card.Footer>
			{/if}
		</Card.Root>
	{/each}
</div>
