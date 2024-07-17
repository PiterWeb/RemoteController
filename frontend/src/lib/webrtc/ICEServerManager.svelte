<script lang="ts">
	export let type: 'stun' | 'turn' = 'stun';

	import DeleteIcon from '$lib/layout/icons/DeleteIcon.svelte';
	import PencilIcon from '$lib/layout/icons/PencilIcon.svelte';
	import TrashIcon from '$lib/layout/icons/TrashIcon.svelte';
	import { ToastType, showToast } from '$lib/toast/toast_hook';

	import autoAnimate from '@formkit/auto-animate';

	import {
		stunServersStore,
		addServerToGroup as addServerToGroupSTUN,
		createServerGroup as createServerGroupSTUN,
		deleteServerGroup as deleteServerGroupSTUN,
		modifyGroup as modifyGroupSTUN,
		removeServerFromGroup as removeServerFromGroupSTUN,
		defaultStunConfig
	} from './stun_servers';

	import {
		turnServersStore,
		addServerToGroup as addServerToGroupTURN,
		createServerGroup as createServerGroupTURN,
		deleteServerGroup as deleteServerGroupTURN,
		modifyGroup as modifyGroupTURN,
		removeServerFromGroup as removeServerFromGroupTURN,
		defaultTurnConfig
	} from './turn_servers';

	import { _ } from 'svelte-i18n';

	const deleteServerGroup = (server_group: string) => {
		if (type === 'stun') {
			deleteServerGroupSTUN(server_group);
		} else {
			deleteServerGroupTURN(server_group);
		}
	};

	const addServerToGroup = (server_group: string) => {
		if (!newserverToAdd) return;
		if (type === 'stun') {
			addServerToGroupSTUN(server_group, newserverToAdd);
		} else {
			addServerToGroupTURN(server_group, newserverToAdd);
		}
	};

	const createServerGroup = () => {
		if (!groupToCreate) return;
		if (type === 'stun') {
			createServerGroupSTUN(groupToCreate);
		} else {
			createServerGroupTURN(groupToCreate);
		}
	};

	const modifyGroup = (
		server_group: string,
		new_group?: string,
		username?: string,
		credential?: string
	) => {
		if (type === 'stun') {
			modifyGroupSTUN(server_group, new_group, username, credential);
		} else {
			modifyGroupTURN(server_group, new_group, username, credential);
		}

		showToast($_('server-group-update'), ToastType.SUCCESS );
	};

	const removeServerFromGroup = (server_group: string, server: string) => {
		if (type === 'stun') {
			removeServerFromGroupSTUN(server_group, server);
		} else {
			removeServerFromGroupTURN(server_group, server);
		}
	};

	const editNameGroup = (server_group: string, i: number) => {
		const el = document.getElementById(`group-${server_group}-${i}`);
		if (!el) return;
		if (el.getAttribute('contenteditable') === 'false') {
			el.setAttribute('contenteditable', 'true');
			el.focus();
			return;
		} else {
			el.setAttribute('contenteditable', 'false');
		}
	};

	let servers = type === 'stun' ? stunServersStore : turnServersStore;

	let groupToCreate: string;
	let newserverToAdd: string;
</script>

<h2 class="text-center text-[clamp(2rem,6vw,4.2rem)] font-black leading-[1.1] xl:text-left">
	<span
		class="[&amp;::selection]:text-base-content text-transparent relative col-start-1 row-start-1 bg-clip-text bg-gradient-to-r from-blue-700 via-blue-800 to-gray-900"
	>
		{#if type === 'stun'}
			{$_('stun-servers-title')}
		{:else}
			{$_('turn-servers-title')}
		{/if}
	</span>
</h2>

<section>
	<button
		on:click={() => {
			if (type === 'stun') {
				const defaultStunConfigCopy = JSON.parse(JSON.stringify(defaultStunConfig));
				stunServersStore.set(defaultStunConfigCopy);
			} else {
				const defaultTurnConfigCopy = JSON.parse(JSON.stringify(defaultTurnConfig));
				turnServersStore.set(defaultTurnConfigCopy);
			}
		}}
		class="btn btn-primary text-white"
	>
		{$_('restore_default_servers')}
	</button>
</section>

<section>
	<form
		on:submit|preventDefault={() => {
			groupToCreate = '';
		}}
		class="flex flex-col gap-4 items-center justify-center sm:w-[30vw] w-[75vw] p-4 my-4 border rounded-lg shadow sm:p-6 bg-gray-800 border-gray-700"
	>
		<label for="group" class="block mb-2 font-medium text-gray-900 dark:text-white"
			>{$_('create_group')}</label
		>
		<input
			type="text"
			id="group"
			class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
			placeholder="Group"
			required
			bind:value={groupToCreate}
		/>
		<button on:click={createServerGroup} type="submit" class="btn btn-primary text-white w-full"
			>{$_('add')}</button
		>
	</form>
</section>

<section id="tutorial-group-server">
	<ul>
		{#if Object.keys($servers).length === 0}
			<p class="text-red-400 text-center text-lg font-medium my-auto h-full">
				{$_('no_groups_warning')}
			</p>
		{/if}
		{#each Object.keys($servers) as server_group, i}
			<li class="w-[75vw] p-4 my-4 border rounded-lg shadow sm:p-6 bg-gray-800 border-gray-700">
				<div class="flex justify-end h-0 mb-4 lg:mb-1">
					<button
						type="button"
						class="btn btn-circle btn-sm btn-ghost text-white"
						on:click={() => deleteServerGroup(server_group)}
					>
						<DeleteIcon /></button
					>
				</div>
				<div class="flex flex-row gap-4">
					<!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
					<!-- svelte-ignore a11y-click-events-have-key-events -->
					<button
						type="button"
						aria-label="Edit group name"
						class="cursor-pointer hover:-rotate-12 transition-transform mb-2"
						on:click={() => editNameGroup(server_group, i)}
					>
						<PencilIcon />
					</button>

					<h4
						class="w-11/12 text-lg font-medium text-gray-900 dark:text-white mb-2 focus"
						contenteditable="false"
						id="group-{server_group}-{i}"
						on:input={(e) =>
							e.currentTarget.textContent && modifyGroup(server_group, e.currentTarget.textContent)}
						on:focusout={() => editNameGroup(server_group, i)}
					>
						{server_group}
					</h4>
				</div>

				<div class="grid lg:grid-cols-3 grid-cols-1 gap-y-10 lg:gap-x-4">
					<form
						on:submit|preventDefault={() => {
							newserverToAdd = '';
						}}
					>
						<label for="domain" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
							>{$_('new_server')}</label
						>
						<input
							type="text"
							id="domain"
							class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg mb-2 focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
							placeholder="domain:port // ip:port"
							required
							bind:value={newserverToAdd}
						/>
						<button
							on:click={() => {
								addServerToGroup(server_group);
							}}
							type="submit"
							class="btn btn-primary text-white w-full">{$_('add')}</button
						>
					</form>

					<ul use:autoAnimate class="max-w-md mx-auto divide-y divide-gray-700 w-full">
						{#if ($servers[server_group]?.urls ?? []).length === 0}
							<div>
								<p class="text-white text-center text-lg font-medium my-auto h-full">
									{$_('no_servers')}
								</p>
								<p class="text-red-400 text-center text-lg font-medium my-auto h-full">
									{$_('no_groups_warning')}
								</p>
							</div>
						{/if}
						{#each $servers[server_group]?.urls ?? [] as server, j}
							<li class="pb-3 sm:pb-4">
								<div class="flex items-center space-x-4 rtl:space-x-reverse">
									<div class="flex-1 min-w-0">
										<p class="text-lg truncate text-white">
											{#if type === 'stun'}
												{server.split('stun:')[1]}
											{:else}
												{server.split('turn:')[1]}
											{/if}
										</p>
									</div>
									<button
										on:click={() => removeServerFromGroup(server_group, server)}
										class="h-5 w-5 mx-10 text-gray-400 cursor-pointer"
									>
										<TrashIcon />
									</button>
								</div>
							</li>
						{/each}
					</ul>

					<div class="flex flex-col justify-center">
						<form action="">
							<label
								for="user-{i}"
								class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
								>{$_('username')}</label
							>
							<input
								type="text"
								id="user-{i}"
								class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
								placeholder="username"
								required
								value={$servers[server_group]?.username ?? ''}
								on:change={(e) => modifyGroup(server_group, undefined, e.currentTarget.value)}
							/>

							<label
								for="password"
								class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
								>{$_('password')}</label
							>

							<input
								type="password"
								id="password"
								class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
								placeholder="•••••••••"
								required
								value={$servers[server_group]?.credential ?? ''}
								on:focusin={(e) => e.currentTarget.type = 'text'}
								on:focusout={(e) => e.currentTarget.type = 'password'}
								on:change={(e) => {
									e.preventDefault()
									const username = $servers[server_group]?.username;
									console.log(username);
									console.log(e.currentTarget.value);
									modifyGroup(server_group, undefined, username, e.currentTarget.value);
								}}
							/>
						</form>
					</div>
				</div>
			</li>
		{/each}
	</ul>
</section>
