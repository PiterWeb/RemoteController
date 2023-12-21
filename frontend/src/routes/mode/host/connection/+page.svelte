<script lang="ts">
	import { startStreaming } from '$lib/webrtc/stream/host_stream_hook';
	import { ListenForConnectionChanges } from '$lib/webrtc/host_webrtc_hook';

	import { onMount } from 'svelte';

	onMount(() => {
		ListenForConnectionChanges();
	});

	interface User {
		name: string;
		ping: number;
	}

	let activeUsers: User[] = [
		// { name: 'John', ping: 50 },
		// { name: 'Jane', ping: 100 },
		// { name: 'Bob', ping: 75 }
	];

	let audio = true;
</script>

<!-- Use TailwindCSS and DaisyUI to style the list -->
<div class="max-w-md mx-auto bg-white rounded-xl shadow-md overflow-hidden md:max-w-2xl">
	<div class="md:flex">
		<div class="p-8">
			<button on:click={() => startStreaming()} class="btn btn-primary">Start Streaming</button
			>
			<ul class="divide-y divide-gray-200">
				{#each activeUsers as user}
					<li class="py-4 flex">
						<div class="ml-3">
							<p class="text-sm font-medium text-gray-900">{user.name}</p>
							<p class="text-sm text-gray-500">{user.ping} ms</p>
						</div>
					</li>
				{/each}
			</ul>
		</div>
	</div>
</div>
