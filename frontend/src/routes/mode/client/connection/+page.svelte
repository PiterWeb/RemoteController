<script lang="ts">
	import { streamingConsumingVideoElement } from '$lib/webrtc/stream/stream_signal_hook';

	let consumingStream = false;

	function connectToStream() {
		streamingConsumingVideoElement.set(
			document.getElementById('stream-video')! as HTMLVideoElement
		);

		consumingStream = true;
	}
</script>

<button class="btn btn-primary" disabled={consumingStream} on:click={connectToStream}>
	Connect to stream</button
>

<video
	id="stream-video"
	class="aspect-video"
	class:hidden={!consumingStream}
	controls
	controlslist="nodownload"
>
	<track kind="captions" />
</video>
