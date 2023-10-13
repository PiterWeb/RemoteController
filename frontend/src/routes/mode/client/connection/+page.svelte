<script lang="ts">
	import { consumeStreaming } from '$lib/hooks/stream';
	import { onMount } from 'svelte';

	let videoElement: HTMLVideoElement;

	onMount(() => {
		consumeStreaming(videoElement);
	});
</script>

<video bind:this={videoElement} />
