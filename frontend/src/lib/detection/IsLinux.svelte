<script lang="ts">
	import { isLinux } from '$lib/detection/detect_os';
	import type { Snippet } from 'svelte';
	let { children, not = false }: {children?: Snippet, not?: boolean} = $props();
</script>

{#await isLinux() then bool}
	{#if bool}
		{#if !not}
			{@render children?.()}
		{/if}
	{:else if not}
		{@render children?.()}
	{/if}
{/await}
