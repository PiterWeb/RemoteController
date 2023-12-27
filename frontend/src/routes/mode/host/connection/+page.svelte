<script lang="ts">
	import { CreateHostStream } from '$lib/webrtc/stream/host_stream_hook';
	import { ListenForConnectionChanges } from '$lib/webrtc/host_webrtc_hook';

	import { onMount } from 'svelte';

	let streaming = false;

	function createStream() {
		CreateHostStream();
		streaming = true;
	}

	onMount(() => {
		ListenForConnectionChanges();
	});
</script>

<button on:click={createStream} disabled={streaming} class="btn btn-primary">Start Streaming</button
>
