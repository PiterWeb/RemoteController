<script>
	import { isWindows } from '$lib/detection/detect_os';
</script>

{#await isWindows() then bool}
	{#if bool}
		<slot />
	{/if}
{/await}
