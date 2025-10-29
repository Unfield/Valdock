<script lang="ts">
	import { CopyCheck, Copy } from '@lucide/svelte/icons';
	import { Button } from '$lib/components/ui/button/index.js';
	import { cn } from '$lib/utils';

	let copied = $state(false);

	let {
		value,
		size = 'icon',
		timeout = 1500,
		className
	}: {
		value: string;
		size?: 'sm' | 'default' | 'icon';
		timeout?: number;
		className?: string;
	} = $props();

	async function handleCopy() {
		try {
			await navigator.clipboard.writeText(value);
			copied = true;
			setTimeout(() => (copied = false), timeout);
		} catch (err) {
			console.error('Copy failed: ', err);
		}
	}
</script>

<Button
	onclick={handleCopy}
	variant="secondary"
	{size}
	class={cn('relative transition-colors hover:text-primary', className)}
	title="Copy to clipboard"
>
	{#if copied}
		<CopyCheck class="h-4 w-4 text-emerald-500 transition-transform" />
	{:else}
		<Copy class="h-4 w-4 text-muted-foreground" />
	{/if}
	<span class="sr-only">Copy to clipboard</span>
</Button>
