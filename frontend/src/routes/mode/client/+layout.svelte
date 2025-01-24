<script>
	import { consumingStream } from '$lib/webrtc/stream/stream_signal_hook.svelte';

    /** @type {{children?: import('svelte').Snippet}} */
    let {children} = $props()

</script>

{@render children?.()}

<video
	id="stream-video"
	class="aspect-video"
	class:hidden={!consumingStream.value}
	controls
	controlslist="nodownload"
>
	<track kind="captions" />
</video>