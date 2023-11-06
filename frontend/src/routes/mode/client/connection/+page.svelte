<script lang="ts">
	import { consumeStreaming } from '$lib/hooks/stream';
	import { onMount } from 'svelte';

	onMount(() => {
		consumeStreaming();
	});
</script>

<video id="stream-video" controls controlslist="nodownload">
	<track kind="captions" />
</video>
