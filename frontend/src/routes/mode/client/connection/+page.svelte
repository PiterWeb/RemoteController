<script lang="ts">
	import { streamingConsumingVideoElement, consumingStream } from '$lib/webrtc/stream/stream_signal_hook.svelte';
	import { _ } from 'svelte-i18n';

	function connectToStream() {
		streamingConsumingVideoElement.set(
			document.getElementById('stream-video')! as HTMLVideoElement
		);

		consumingStream.value = true;
	}
</script>

<button class="btn btn-primary" disabled={consumingStream.value} onclick={connectToStream}>
	{$_('connect-to-stream')}</button
>

<video
	id="stream-video"
	class="aspect-video"
	class:hidden={!consumingStream.value}
	controls
	controlslist="nodownload"
>
	<track kind="captions" />
</video>
