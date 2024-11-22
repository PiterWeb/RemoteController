<script>
	import { isWindows } from '$lib/detection/detect_os';
	/** @type {{children?: import('svelte').Snippet}} */
	let { children } = $props();
</script>

{#await isWindows() then bool}
	{#if bool}
		{@render children?.()}
	{/if}
{/await}
