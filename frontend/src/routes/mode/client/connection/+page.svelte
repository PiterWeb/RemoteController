<script lang="ts">
	import { CloseStreamClientConnection } from '$lib/webrtc/stream/client_stream_hook';
	import {consumingStream } from '$lib/webrtc/stream/stream_signal_hook.svelte';
	import { _ } from 'svelte-i18n';

	function connectToStream() {
		consumingStream.value = true;
	}

	function closeStream() {
		CloseStreamClientConnection()
	}

</script>


	<button class="btn btn-primary" class:hidden={consumingStream.value} onclick={connectToStream}>
		{$_('connect-to-stream')}</button>

<button class="btn btn-primary btn-outline" class:hidden={!consumingStream.value} onclick={closeStream}>
	{$_('close-stream')}</button>