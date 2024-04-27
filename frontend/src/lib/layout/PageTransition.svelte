<script lang="ts">
	import { slide } from 'svelte/transition';

	export let key: string;
	export let duration: number = 0;

	const audio = new Audio('/sounds/page_transition.mp3');

	audio.volume = 0.1;

	$: key && navigator.userActivation.isActive && audio.play();

</script>

{#key key}
	<div in:slide={{ delay: duration / 4, duration, axis: 'x' }}>
		<slot />
	</div>
{/key}
