<script lang="ts">

    import { consumingStream, mediaStreams } from '$lib/webrtc/stream/stream_signal_hook.svelte';

    $inspect(consumingStream);

	let { children } = $props();

	function srcObject(node: HTMLVideoElement, stream: MediaStream) {
		node.srcObject = stream;
		return {
			update(nextStream: MediaStream) { node.srcObject = stream;  },
			destroy() { /* stream revoking logic here */ },
		}
	}

</script>

{@render children?.()}

<div
	id="stream-video"
	class="w-full h-full"
	class:hidden={!consumingStream.value}
>	

	{#each mediaStreams.value as stream}
		<video use:srcObject={stream} autoplay controls>
			<track kind="captions"/>
		</video>
	{/each}

</div>