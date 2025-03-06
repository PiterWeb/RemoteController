<script lang="ts">
	import { CreateHostStream } from '$lib/webrtc/stream/host_stream_hook';
	import { DEFAULT_IDEAL_FRAMERATE, DEFAULT_MAX_FRAMERATE, FIXED_RESOLUTIONS } from '$lib/webrtc/stream/stream_config';
	import { _ } from 'svelte-i18n';
	import { streaming } from '$lib/webrtc/stream/stream_signal_hook.svelte';
	import { ListenForConnectionChanges } from '$lib/webrtc/host_webrtc_hook';

	import { onMount } from 'svelte';
	import IsLinux from '$lib/detection/IsLinux.svelte';

	let selected_resolution = $state(FIXED_RESOLUTIONS.resolution720p);

	let idealFramerate = $state(DEFAULT_IDEAL_FRAMERATE);
	let maxFramerate = $state(DEFAULT_MAX_FRAMERATE);

	function createStream() {
		CreateHostStream(selected_resolution, idealFramerate, maxFramerate);
		streaming.value = true;
	}

	onMount(() => {
		ListenForConnectionChanges();
	});
</script>

<IsLinux>
	<div class="w-full h-full">
		<h3 class="text-4xl">{$_('relay-title')}</h3>
		<p class="text-lg">{$_('go-browser')}</p>
		<p class="text-error">{$_('warning-go-browser')}</p>
	</div>
</IsLinux>

<IsLinux not>
	<div class:hidden={streaming.value} class="w-full h-full grid grid-rows-2 grid-cols-1 gap-2">
		<div class="w-full h-full">
			<h3 class="text-3xl">{$_('resolutions')}</h3>
			<select
				class="select select-primary w-full max-w-xs mt-6"
				bind:value={selected_resolution}
				id="resolution"
				aria-label="resolution"
			>
				{#each Object.values(FIXED_RESOLUTIONS) as resolution}
					<option selected={resolution === selected_resolution} value={resolution}
						>{resolution}p</option
					>
				{/each}
			</select>
		</div>
		<div class="w-full h-full">
			<h3 class="text-3xl">{$_('framerate')}</h3>
			<div class="flex flex-row gap-10 h-full w-full mt-6">
				<div class="w-full h-full flex flex-col">
					<h4 class="text-lg">{$_('ideal-framerate')}</h4>
						<input
						type="number"
						class="w-10 h-10 bg-neutral text-white text-center"
						bind:value={idealFramerate}
						pattern="0|[1-9]\d*"
						min="25"
						max="145"
						step="5"
						/>
					<input
						type="range"
						min="25"
						max="145"
						bind:value={idealFramerate}
						class="range range-lg my-10"
						step="5"
					/>
				</div>
				<div class="w-full h-full flex flex-col">
					<h4 class="text-lg">{$_('max-framerate')}</h4>
					<input
					type="number"
					class="w-10 h-10 bg-neutral text-white text-center"
					bind:value={maxFramerate}
					pattern="0|[1-9]\d*"
					min="25"
					max="145"
					step="5"
					/>
					<input
						type="range"
						min="30"
						max="145"
						bind:value={maxFramerate}
						class="range range-lg my-10"
						step="5"
					/>
				</div>
			</div>
		</div>
	</div>

	<button onclick={createStream} disabled={streaming.value} class="btn btn-primary"
		>{$_('start-streaming')}</button
	>
</IsLinux>

<style>

	input[type="number"] {
		-webkit-appearance: textfield;
		-moz-appearance: textfield;
		appearance: textfield;
	}

</style>
