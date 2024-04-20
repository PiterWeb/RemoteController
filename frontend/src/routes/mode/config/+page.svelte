<script script lang="ts">
	import StunServers from './StunServers.svelte';
	import IsWindows from '$lib/detection/IsWindows.svelte';
	import ViGEmDownload from './ViGEmDownload.svelte';
	import { _ } from 'svelte-i18n';
	import LanguageSelector from './LanguageSelector.svelte';
	import TurnServers from './TurnServers.svelte';
</script>

<h2 class="text-center text-[clamp(2rem,6vw,4.2rem)] font-black leading-[1.1] xl:text-left">
	<span
		class="[&amp;::selection]:text-base-content text-transparent relative col-start-1 row-start-1 bg-clip-text bg-gradient-to-r from-blue-700 via-blue-800 to-gray-900"
		>{$_('config_title')}
	</span>
</h2>

<div class="grid gap-28 w-full md:grid-cols-3 py-12 xl:[&>*]:w-80 lg:[&>*]:w-72 md:[&>*]:w-60 [&>*]:w-72">
	<LanguageSelector />
	<StunServers />
	<TurnServers />

	<IsWindows>
		<ViGEmDownload />
	</IsWindows>
</div>
