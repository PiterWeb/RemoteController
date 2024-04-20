<script lang="ts">
	import { slide } from 'svelte/transition';

	export let key: string;
	export let duration: number = 0;
</script>

{#key key}
	<div in:slide={{ delay: duration / 4, duration, axis: 'x' }}>
		<slot />
	</div>
{/key}
