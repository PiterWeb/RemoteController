<script lang="ts">
	import { _ } from 'svelte-i18n';
	import { consumingStream } from '$lib/webrtc/stream/stream_signal_hook.svelte';

	function toogleStream() {
		consumingStream.value = !consumingStream.value;
	}
</script>

<button class="btn btn-primary" class:btn-neutral={consumingStream.value}  onclick={toogleStream}>
	{$_('connect-to-stream')}</button
>


