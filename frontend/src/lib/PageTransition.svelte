<script lang="ts">
	import { slide } from 'svelte/transition';

	export let key: string;
	export let duration: number = 0;

	const isReduced =
		window.matchMedia(`(prefers-reduced-motion: reduce)`).matches === true ||
		window.matchMedia(`(prefers-reduced-motion: reduce)`).matches === true;
</script>

{#if isReduced}
	<slot />
{:else}
	{#key key}
		<div in:slide={{ delay: duration / 4, duration, axis: 'x' }}>
			<slot />
		</div>
	{/key}
{/if}
