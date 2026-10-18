<script lang="ts">
	import { CloseStreamClientConnection } from '$lib/webrtc/stream/client_stream_hook';
	import { videoSpeedOptimizationEnabled } from '$lib/webrtc/stream/stream_config.svelte';
	import {consumingStream } from '$lib/webrtc/stream/stream_signal_hook.svelte';
	import { _ } from 'svelte-i18n';
	
	function connectToStream() {
		consumingStream.value = true;
	}

	function closeStream() {
		CloseStreamClientConnection()
	}
	
</script>

<div class="grid md:grid-cols-4 grid-cols-1 gap-10">
	{#if consumingStream.value}
	<button class="md:col-span-4 btn btn-primary btn-outline" onclick={closeStream}>
		{$_('close-stream')}
	</button>
	
	{:else}

	<section class="flex flex-row-reverse items-center gap-2">
		<label
			id="label-external-whip-checkbox"
			for="lan-mode-checkbox"
			class="font-semibold text-white">{$_('config-video-speed-optimization')}</label
		>
		<input
			id="lan-mode-checkbox"
			type="checkbox"
			class="checkbox checkbox-xs checkbox-primary"
			checked={videoSpeedOptimizationEnabled.value}
			onchange={() => videoSpeedOptimizationEnabled.value = !videoSpeedOptimizationEnabled.value}
		/>
	</section>
	
	<button class="md:col-span-2 btn btn-primary" onclick={connectToStream}>
		{$_('connect-to-stream')}</button>
		
	{/if}
</div>

