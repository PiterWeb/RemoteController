<script lang="ts">
	import { slide } from 'svelte/transition';
	import { playAudio } from '$lib/audio/audio_player';

	interface Props {
		key: string;
		duration?: number;
		children?: import('svelte').Snippet;
	}

	let { key, duration = 0, children }: Props = $props();

	$effect(() => {
		key && playAudio('page_transition');
	});
</script>

{#key key}
	<div in:slide={{ delay: duration / 4, duration, axis: 'x' }}>
		{@render children?.()}
	</div>
{/key}
