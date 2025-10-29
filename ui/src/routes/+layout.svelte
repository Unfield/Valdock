<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import { ModeWatcher } from 'mode-watcher';

	let { children } = $props();
</script>

<svelte:head>
	<title>Valdock WebUI</title>
	<link rel="icon" href={favicon} />
</svelte:head>

<ModeWatcher defaultMode="system" />
{@render children?.()}
