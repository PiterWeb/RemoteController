<script lang="ts">
	import { CreateHostStream, StopStreaming } from '$lib/webrtc/stream/host_stream_hook';
	import {
		DEFAULT_IDEAL_FRAMERATE,
		DEFAULT_MAX_FRAMERATE,
		FIXED_RESOLUTIONS,
		MAX_FRAMES,
		MIN_FRAMES,
		RESOLUTIONS
	} from '$lib/webrtc/stream/stream_config.svelte';
	import { _ } from 'svelte-i18n';
	import { streaming } from '$lib/webrtc/stream/stream_signal_hook.svelte';
	import { ListenForConnectionChanges } from '$lib/webrtc/host_webrtc_hook';

	import { onMount } from 'svelte';
	import CodecList from '$lib/webrtc/stream/CodecList.svelte';
	import IsLinux from '$lib/detection/IsLinux.svelte';
	import KeyboardIcon from '$lib/layout/icons/KeyboardIcon.svelte';
	import GamepadIcon from '$lib/layout/icons/GamepadIcon.svelte';
	import type { audio } from '$lib/wailsjs/go/models';
	import ws from '$lib/websocket/ws';
	import { IS_RUNNING_EXTERNAL } from '$lib/detection/onwebsite';
	import Tooltip from '$lib/layout/Tooltip.svelte';
	import MouseIcon from '$lib/layout/icons/MouseIcon.svelte';
	import IsRunningExternal from '$lib/detection/IsRunningExternal.svelte';
	// import { GetAudioProcess, SetAudioPid } from '$lib/bindings/App';

	let selected_audio_src = $state(0);
	let audio_srcs = $state<audio.AudioProcess[]>([]);

	let selected_resolution = $state(FIXED_RESOLUTIONS.resolution720p);

	let idealFramerate = $state(DEFAULT_IDEAL_FRAMERATE);
	let maxFramerate = $state(DEFAULT_MAX_FRAMERATE);

	let canStartStreaming = $state(false);

	let timeoutSetAudioPid: NodeJS.Timeout;

	$effect(() => {
		// *** This is a hack needed to trigger the effect
		selected_audio_src;
		// ***
		if (streaming.value) {
			clearTimeout(timeoutSetAudioPid);
			timeoutSetAudioPid = setTimeout(() => {
				console.log('Selected audio :', selected_audio_src);
				// SetAudioPid(selected_audio_src)
			}, 750);
		} else {
			// SetAudioPid(0)
		}
	});

	let whipEnabled = $state(false);

	$effect(() => {
		if (whipEnabled) streaming.value = true;
		else streaming.value = false;
	});

	$effect(() => {
		if (idealFramerate > maxFramerate) idealFramerate = maxFramerate;
		else if (idealFramerate < MIN_FRAMES) idealFramerate = MIN_FRAMES;
	});

	$effect(() => {
		if (maxFramerate > MAX_FRAMES) maxFramerate = MAX_FRAMES;
		else if (maxFramerate < MIN_FRAMES) maxFramerate = MIN_FRAMES;
	});

	let keyboardEnabled = $state(false);
	let gamepadEnabled = $state(true);
	let mouseEnabled = $state(false);

	let httpExternalPort = $state(8090);
	let whipPort = $state(8092);

	function createStream() {
		CreateHostStream(selected_resolution, idealFramerate, maxFramerate);
		streaming.value = true;
		canStartStreaming = false;
	}

	async function toogleWhip() {
		const { IsWhipEnabled, ToogleWhip } = await import('$lib/bindings/App');

		await ToogleWhip();
		whipEnabled = await IsWhipEnabled();
	}

	async function toogleKeyboard() {
		const { IsKeyboardEnabled, ToogleKeyboard } = await import('$lib/bindings/App');

		await ToogleKeyboard();
		keyboardEnabled = await IsKeyboardEnabled();
	}

	async function toogleGamepad() {
		const { IsGamepadEnabled, ToogleGamepad } = await import('$lib/bindings/App');

		await ToogleGamepad();
		gamepadEnabled = await IsGamepadEnabled();
	}

	async function toogleMouse() {
		const { IsMouseEnabled, ToogleMouse } = await import('$lib/bindings/App');
		await ToogleMouse();
		mouseEnabled = await IsMouseEnabled();
	}

	function MapAudioSrcs(s: audio.AudioProcess) {
		if (s.Name.length > 0) return s;
		else return { Name: `Unknow<${s.Pid}>`, Pid: s.Pid };
	}

	onMount(() => {
		ListenForConnectionChanges();

		(async () => {
			// audio_srcs = (await GetAudioProcess()).map(MapAudioSrcs)
		})();

		const interval = setInterval(async () => {
			// audio_srcs = (await GetAudioProcess()).map(MapAudioSrcs)
		}, 5000);

		let unlistener: () => void;

		// Handle start streaming button
		(async () => {
			if (IS_RUNNING_EXTERNAL) {
				const cllbck = (ev: MessageEvent<string>) => {
					if (ev.data == '{}') canStartStreaming = true;
				};
				ws().addEventListener('message', cllbck);
				unlistener = () => ws().removeEventListener('message', cllbck);
				return;
			}

			const { IsWhipEnabled } = await import('$lib/bindings/App');
			whipEnabled = await IsWhipEnabled();

			const { EventsOn } = await import('$lib/wailsjs/runtime/runtime');
			const { GetUsedPorts } = await import('$lib/bindings/App');
			const usedPorts = await GetUsedPorts();
			httpExternalPort = usedPorts.HTTP;
			whipPort = usedPorts.WHIP;
			unlistener = EventsOn('streaming-signal-client', (d) => {
				if (d == '{}') canStartStreaming = true;
			});
		})();

		return () => {
			unlistener();
			clearInterval(interval);
		};
	});
</script>

<Tooltip ref="label-external-whip-checkbox" placement="top">
	<p>{$_('external-whip-explanation')}</p>
</Tooltip>

<IsRunningExternal not>
	<section
		class:hidden={streaming.value && !whipEnabled}
		class="flex flex-row-reverse items-center gap-2"
	>
		<label
			id="label-external-whip-checkbox"
			for="lan-mode-checkbox"
			class="font-semibold text-white">{$_('external-whip')}</label
		>
		<input
			id="lan-mode-checkbox"
			type="checkbox"
			class="checkbox checkbox-xs checkbox-primary"
			checked={whipEnabled}
			onchange={toogleWhip}
		/>
	</section>

	<section class:hidden={!whipEnabled} class="text-white">
		Whip URL: http://localhost:{whipPort}/whip
	</section>

	<section class="w-full">
		<h3 class="text-3xl text-white text-center">
			{$_('toogle_devices')}
		</h3>
		<div class="flex flex-row justify-center gap-3 mt-6">
			<button
				onclick={toogleMouse}
				class:btn-primary={mouseEnabled}
				class:btn-neutral={!mouseEnabled}
				class:border-gray-400={!mouseEnabled}
				class="btn border"
			>
				<MouseIcon />
			</button>
			<button
				onclick={toogleKeyboard}
				class:btn-primary={keyboardEnabled}
				class:btn-neutral={!keyboardEnabled}
				class:border-gray-400={!keyboardEnabled}
				class="btn border"
			>
				<KeyboardIcon />
			</button>
			<button
				onclick={toogleGamepad}
				class:btn-primary={gamepadEnabled}
				class:btn-neutral={!gamepadEnabled}
				class:border-gray-400={!gamepadEnabled}
				class="btn border"
			>
				<GamepadIcon />
			</button>
		</div>
		<div
			class:hidden={!streaming.value || whipEnabled}
			class="flex flex-row justify-center gap-3 mt-6"
		>
			<button onclick={StopStreaming} disabled={!streaming.value} class="btn btn-primary">
				{$_('stop-streaming')}
			</button>
		</div>
	</section>
</IsRunningExternal>

<!-- <section class="w-1/3 mx-auto" class:hidden={!streaming.value}>
	<h3 class="text-3xl text-white text-center">
		{$_('audio_selector')}
	</h3>
	<select
			class="select w-full mx-auto mt-6"
			bind:value={selected_audio_src}
			id="audio_srcs"
			aria-label="audio_srcs"
		>
			{#each audio_srcs as src}
				<option selected={src.Pid == selected_audio_src} value={src.Pid}
					>{src.Name}</option
				>
			{/each}
		</select>
</section> -->

<IsLinux>
	<div class="w-full h-full">
		<h3 class="text-4xl text-white">{$_('relay-title')}</h3>
		<p class="text-gray-300">http://localhost:{httpExternalPort}/mode/host/connection/</p>
		<p class="text-lg text-gray-400">{$_('go-browser')}</p>
		<p class="text-error">{$_('warning-go-browser')}</p>
	</div>
</IsLinux>

<IsLinux not>
	<div class:hidden={streaming.value} class="w-2/3 flex flex-col md:flex-row gap-12 align-middle">
		<CodecList />
		<section class="w-full">
			<h3 class="text-3xl text-white text-center">{$_('resolutions')}</h3>
			<select
				class="select w-full mx-auto mt-6"
				bind:value={selected_resolution}
				id="resolution"
				aria-label="resolution"
			>
				{#each Object.values(FIXED_RESOLUTIONS) as resolution}
					<option selected={resolution === selected_resolution} value={resolution}>
						{resolution}
						{#if resolution !== FIXED_RESOLUTIONS.resolutionNative}
							p
						{/if}
					</option>
				{/each}
			</select>
		</section>
		<section class="w-full">
			<h3 class="text-3xl text-white text-center">{$_('framerate')}</h3>
			<div class="flex flex-row gap-10 w-full mt-6">
				<div class="w-full flex flex-col">
					<h4 class="text-lg text-gray-300">{$_('ideal-framerate')}</h4>
					<input
						type="number"
						class="input w-16 px-2 h-10 bg-neutral border border-gray-400 text-white text-center"
						bind:value={idealFramerate}
						pattern="0|[1-9]\d*"
						min={MIN_FRAMES}
						max={maxFramerate}
						step="5"
					/>
					<input
						type="range"
						min={MIN_FRAMES}
						max={maxFramerate}
						bind:value={idealFramerate}
						class="range range-primary bg-white range-lg my-10"
						step="5"
					/>
				</div>
				<div class="w-full flex flex-col">
					<h4 class="text-lg text-gray-300">{$_('max-framerate')}</h4>
					<input
						type="number"
						class="input w-16 px-2 h-10 bg-neutral border border-gray-400 text-white text-center"
						bind:value={maxFramerate}
						pattern="0|[1-9]\d*"
						min={MIN_FRAMES}
						max={MAX_FRAMES}
						step="5"
					/>
					<input
						type="range"
						min={MIN_FRAMES}
						max={MAX_FRAMES}
						bind:value={maxFramerate}
						class="range range-primary bg-white range-lg my-10"
						step="5"
					/>
				</div>
			</div>
		</section>
	</div>

	<button
		onclick={createStream}
		disabled={streaming.value || !canStartStreaming}
		class="btn btn-primary">{$_('start-streaming')}</button
	>
</IsLinux>

<style>
	input[type='number'] {
		-webkit-appearance: textfield;
		-moz-appearance: textfield;
		appearance: textfield;
	}
</style>
