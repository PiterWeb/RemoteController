<script lang="ts">
	import { consumeStreaming } from '$lib/hooks/client_stream';
	
	let consuming = false;

	function connectToStream() {
		consuming = true;
		consumeStreaming();
	}

</script>

<button class="btn btn-primary" disabled={consuming} on:click={connectToStream}> Connect to stream</button>

<video id="stream-video" controls controlslist="nodownload">
	<track kind="captions" />
</video>
