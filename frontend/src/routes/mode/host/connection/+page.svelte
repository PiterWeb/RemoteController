<script lang="ts">
	import { CreateHostStream } from '$lib/webrtc/stream/host_stream_hook';
	import { ListenForConnectionChanges } from '$lib/webrtc/host_webrtc_hook';

	import { onMount } from 'svelte';

	let streaming = false

	function createStream() {
		CreateHostStream();
		streaming = true;
	}

	onMount(() => {
		ListenForConnectionChanges();
	});

</script>

<!-- Use TailwindCSS and DaisyUI to style the list -->
<div class="max-w-md mx-auto bg-white rounded-xl shadow-md overflow-hidden md:max-w-2xl">
	<div class="md:flex">
		<div class="p-8">
			<button on:click={createStream} disabled={streaming} class="btn btn-primary">Start Streaming</button
			>
		</div>
	</div>
</div>
