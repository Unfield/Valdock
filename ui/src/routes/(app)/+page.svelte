<script lang="ts">
	import Loading from '$lib/components/loading.svelte';
	import { onMount } from 'svelte';
	import Cardview from './components/cardview.svelte';

	let loading = $state(true);

	let instances = $state<App.Instance[]>([]);

	async function fetchInstances() {
		loading = true;
		const res = await fetch('/api/v1/instances');
		const json = await res.json();

		instances = json.data.instances.map((i: App.Instance) => ({
			id: i.id,
			name: i.name,
			port: i.port,
			status: i.status,
			createdAt: i.createdAt
		}));
		loading = false;
	}

	onMount(() => {
		fetchInstances();
	});
</script>

<svelte:head>
	<title>Dashboard - Valdock WebUI</title>
</svelte:head>

{#if loading && instances.length < 1}
	<Loading />
{:else}
	<Cardview {instances} />
{/if}
