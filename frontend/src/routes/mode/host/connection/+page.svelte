<script lang="ts">
	import { CreateHostStream } from '$lib/webrtc/stream/host_stream_hook';
	import { FIXED_RESOLUTIONS } from '$lib/webrtc/stream/stream_config';
	import { _ } from 'svelte-i18n';
	import { streaming } from '$lib/webrtc/stream/stream_signal_hook.svelte';
	import { slide } from 'svelte/transition';
	import { ListenForConnectionChanges } from "$lib/webrtc/host_webrtc_hook"

	import { onMount } from 'svelte';
	import { elasticOut } from 'svelte/easing';
	import IsLinux from '$lib/detection/IsLinux.svelte';

	let selected_resolution = $state(FIXED_RESOLUTIONS.resolution720p);

	let idealFramerate = $state(25);

	function createStream() {
		CreateHostStream(selected_resolution);
		streaming.value = true;
	}

	onMount(() => {
		ListenForConnectionChanges();
	});
</script>

<IsLinux>
	<div class="w-full h-full">
		<h3>{$_('relay-title')}</h3>
		<p>{$_('go-browser')}</p>
	</div>
</IsLinux>

<IsLinux not>
	<div class:hidden={streaming.value} class="w-full h-full grid grid-rows-2 grid-cols-1 gap-10">
		<div class="w-full h-full">
			<h3>{$_('resolutions')}</h3>
			<select
				class="select select-primary w-full max-w-xs mt-6"
				bind:value={selected_resolution}
				id="resolution"
				aria-label="Default select example"
			>
				{#each Object.values(FIXED_RESOLUTIONS) as resolution}
					<option selected={resolution === selected_resolution} value={resolution}
						>{resolution}p</option
					>
				{/each}
			</select>
		</div>
		<div class="w-full h-full">
			<h3 class="">{$_('framerate')}</h3>
			<div class="flex flex-row gap-10 h-full w-full mt-6">
				<div class="w-full h-full flex flex-col">
					<h4>{$_('ideal-framerate')}</h4>
					<div class="w-10 h-10 bg-neutral">
						{#key idealFramerate}
							<p transition:slide={{ easing: elasticOut }} class="text-white text-center">
								{idealFramerate}
							</p>
						{/key}
					</div>
					<input
						type="range"
						min="25"
						max="145"
						bind:value={idealFramerate}
						class="range range-lg"
						step="5"
					/>
				</div>
				<div class="w-full h-full flex flex-col">
					<h4>{$_('max-framerate')}</h4>
					<input type="range" min="30" max="150" value="30" class="range range-lg" step="30" />
					<div class="flex w-full justify-between px-2 text-xs">
						<span>30</span>
						<span>60</span>
						<span>90</span>
						<span>120</span>
						<span>150</span>
					</div>
				</div>
			</div>
		</div>
	</div>

	<button onclick={createStream} disabled={streaming.value} class="btn btn-primary"
		>{$_('start-streaming')}</button
	>
</IsLinux>
