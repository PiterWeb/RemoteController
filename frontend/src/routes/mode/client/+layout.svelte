<script>

    import { consumingStream } from '$lib/webrtc/stream/stream_signal_hook.svelte';

    $inspect(consumingStream);

	let { children } = $props();
</script>

{@render children?.()}

<div
	id="stream-video"
	class="w-full h-full"
	class:hidden={!consumingStream.value}
>	
</div>