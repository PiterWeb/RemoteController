<script lang="ts">

    import turnServers from '$lib/webrtc/turn_servers';
    import {_} from 'svelte-i18n';

    function removeTurnServer(index: number) {
		turnServers.update((servers) => {
			servers.splice(index, 1);
			return servers;
		});
	}

	function submitForm(e: SubmitEvent) {
		const form = e.target as HTMLFormElement;

		const formData = new FormData(form);

		const stunServer = (formData.get('turn-server') ?? '') as string;

		if (stunServer.length === 0) return;

		turnServers.update((servers) => [...servers, 'turn:' + stunServer]);

		form.reset();
	}

</script>

<section
		class="max-w-96 min-w-full p-4 bg-white border border-gray-200 rounded-lg shadow sm:p-6 md:p-8 dark:bg-gray-800 dark:border-gray-700"
	>
		<form class="max-w-md mx-auto space-y-6" on:submit|preventDefault={submitForm}>
			<label for="helper-text" class="text-xl font-medium text-gray-900 dark:text-white"
				>{$_('turn-servers-title')}</label
			>

            <a
				href="https://github.com/coturn/coturn"
				target="_blank"
				class="text-sm font-medium text-blue-600 hover:underline dark:text-blue-500"
				>{$_('selfhost-turn-link')}</a
			>

			<div class="relative">
				<input
					type="text"
					id="helper-text"
					name="stun-server"
					aria-describedby="helper-text-explanation"
					class="block w-full p-4 text-sm text-gray-900 border border-gray-300 rounded-lg bg-gray-50 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
					placeholder="turn.example.com:3478"
				/>
				<button
					type="submit"
					class="text-white absolute end-2.5 bottom-2.5 bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-4 py-2 dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
					>{$_('add')}</button
				>
			</div>

			<ul class="max-w-lg space-y-4 text-gray-500 list-inside dark:text-gray-400">
				{#each $turnServers as server, i}
					<li class="flex items-center">
						<svg
							class="w-3.5 h-3.5 me-2 text-green-500 dark:text-green-400 flex-shrink-0"
							aria-hidden="true"
							xmlns="http://www.w3.org/2000/svg"
							fill="currentColor"
							viewBox="0 0 20 20"
						>
							<path
								d="M10 .5a9.5 9.5 0 1 0 9.5 9.5A9.51 9.51 0 0 0 10 .5Zm3.707 8.207-4 4a1 1 0 0 1-1.414 0l-2-2a1 1 0 0 1 1.414-1.414L9 10.586l3.293-3.293a1 1 0 0 1 1.414 1.414Z"
							/>
						</svg>
						{server}
						<button
							on:click={() => removeTurnServer(i)}
							type="button"
							class="ms-auto -mx-1.5 -my-1.5 bg-white text-gray-400 hover:text-gray-900 rounded-lg focus:ring-2 focus:ring-gray-300 p-1.5 hover:bg-gray-100 inline-flex items-center justify-center h-8 w-8 dark:text-gray-500 dark:hover:text-white dark:bg-gray-800 dark:hover:bg-gray-700"
							data-dismiss-target="#toast-default"
							aria-label="Close"
						>
							<span class="sr-only">Close</span>
							<svg
								class="w-3 h-3"
								aria-hidden="true"
								xmlns="http://www.w3.org/2000/svg"
								fill="none"
								viewBox="0 0 14 14"
							>
								<path
									stroke="currentColor"
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="m1 1 6 6m0 0 6 6M7 7l6-6M7 7l-6 6"
								/>
							</svg>
						</button>
					</li>
				{/each}
			</ul>
		</form>
	</section>