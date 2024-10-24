<script lang="ts">
	import { CreateHostStream } from '$lib/webrtc/stream/host_stream_hook';
	import { ListenForConnectionChanges } from '$lib/webrtc/host_webrtc_hook';
	import { _ } from 'svelte-i18n'

	import { onMount } from 'svelte';

	let streaming = $state(false);

	function createStream() {
		CreateHostStream();
		streaming = true;
	}

	onMount(() => {
		ListenForConnectionChanges();
	});
</script>

<button onclick={createStream} disabled={streaming} class="btn btn-primary">{$_('start-streaming')}</button
>
