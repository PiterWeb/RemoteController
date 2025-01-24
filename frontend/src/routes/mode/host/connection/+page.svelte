<script lang="ts">
	import { CreateHostStream, fixedResolutions } from '$lib/webrtc/stream/host_stream_hook';
	import { ListenForConnectionChanges } from '$lib/webrtc/host_webrtc_hook';
	import { _ } from 'svelte-i18n'
	import { streaming } from "$lib/webrtc/stream/stream_signal_hook.svelte";

	import { onMount } from 'svelte';

	let selected_resolution = $state(fixedResolutions.resolution720p)

	function createStream() {
		CreateHostStream(selected_resolution);
		streaming.value = true;
	}

	onMount(() => {
		ListenForConnectionChanges();
	});
</script>

<div class:hidden={streaming.value} class="w-full h-full">
	<h2>{$_("resolutions")}</h2>
	<select class="select select-primary w-full max-w-xs"  bind:value={selected_resolution} id="resolution" aria-label="Default select example">
		{#each Object.values(fixedResolutions) as resolution}
			<option selected={resolution === selected_resolution} value={resolution}>{resolution}</option>
		{/each}
	</select>
</div>

<button onclick={createStream} disabled={streaming.value} class="btn btn-primary">{$_('start-streaming')}</button>