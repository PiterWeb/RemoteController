<script lang="ts">
	import { slide } from 'svelte/transition';
	import { playAudio } from '$lib/audio/audio_player';

	export let key: string;
	export let duration: number = 0;

	$: key && playAudio('page_transition');
</script>

{#key key}
	<div in:slide={{ delay: duration / 4, duration, axis: 'x' }}>
		<slot />
	</div>
{/key}
