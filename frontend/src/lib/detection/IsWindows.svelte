<script lang="ts">
	import { isWindows } from '$lib/detection/detect_os';
	import type { Snippet } from 'svelte';
	let { children, not = false }: {children?: Snippet, not?: boolean} = $props();
</script>

{#await isWindows() then bool}
	{#if bool}
		{#if !not}
			{@render children?.()}
		{/if}
	{:else if not}
		{@render children?.()}
	{/if}
{/await}
