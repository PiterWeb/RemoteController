<script lang="ts">
	import { streamingConsuming } from "$lib/webrtc/stream/stream_signal_hook";

	let streaming = false;

	function connectToStream() {
		streamingConsuming.set(document.getElementById("stream-video")! as HTMLVideoElement);
		streaming = true;
	}
</script>

<button class="btn btn-primary" disabled={streaming} on:click={connectToStream}>
	Connect to stream</button
>

<video id="stream-video" class="aspect-video" controls controlslist="nodownload">
		<track kind="captions" />
</video>
