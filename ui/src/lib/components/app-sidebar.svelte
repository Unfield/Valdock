<script lang="ts" module>
	import {
		Database,
		User,
		NotepadText,
		Settings,
		ChartArea,
		DatabaseBackup,
		Waypoints
	} from '@lucide/svelte/icons';

	const data = {
		user: {
			name: 'shadcn',
			email: 'm@example.com',
			avatar: '/avatars/shadcn.jpg'
		},
		teams: [
			{
				name: 'Personal',
				logo: User,
				plan: 'Peronal'
			}
		],
		navMain: [
			{
				title: 'Instances',
				url: '#',
				icon: Database,
				isActive: true,
				items: [
					{
						title: 'Overview',
						url: '#'
					},
					{
						title: 'Create',
						url: '#'
					}
				]
			},
			{
				title: 'Templates',
				url: '#',
				icon: NotepadText,
				items: [
					{
						title: 'Overview',
						url: '#'
					},
					{
						title: 'Create',
						url: '#'
					}
				]
			},
			{
				title: 'Clusters',
				url: '#',
				icon: Waypoints,
				items: [
					{
						title: 'Overview',
						url: '#'
					},
					{
						title: 'Create',
						url: '#'
					}
				]
			},
			{
				title: 'Backups',
				url: '#',
				icon: DatabaseBackup,
				items: [
					{
						title: 'Overview',
						url: '#'
					},
					{
						title: 'Create',
						url: '#'
					},
					{
						title: 'Quick Restore',
						url: '#'
					}
				]
			},
			{
				title: 'Metrics',
				url: '#',
				icon: ChartArea,
				items: []
			},
			{
				title: 'Settings',
				url: '#',
				icon: Settings,
				items: [
					{
						title: 'Host',
						url: '#'
					}
				]
			}
		]
	};
</script>

<script lang="ts">
	import NavMain from './nav-main.svelte';
	import NavUser from './nav-user.svelte';
	import TeamSwitcher from './team-switcher.svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import type { ComponentProps } from 'svelte';
	let {
		ref = $bindable(null),
		collapsible = 'icon',
		...restProps
	}: ComponentProps<typeof Sidebar.Root> = $props();
</script>

<Sidebar.Root {collapsible} {...restProps}>
	<Sidebar.Header>
		<TeamSwitcher teams={data.teams} />
	</Sidebar.Header>
	<Sidebar.Content>
		<NavMain items={data.navMain} />
	</Sidebar.Content>
	<Sidebar.Footer>
		<NavUser user={data.user} />
	</Sidebar.Footer>
	<Sidebar.Rail />
</Sidebar.Root>
