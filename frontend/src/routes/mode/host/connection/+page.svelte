<script lang="ts">
	import { CreateHostStream } from '$lib/webrtc/stream/host_stream_hook';
	import { ListenForConnectionChanges } from '$lib/webrtc/host_webrtc_hook';
	import { _ } from 'svelte-i18n'

	import { onMount } from 'svelte';

	onMount(() => {
		ListenForConnectionChanges();
	});
</script>

<h2>Connection</h2>