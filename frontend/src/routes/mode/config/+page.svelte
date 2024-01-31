<script script lang="ts">
	import StunServers from '../../../lib/webrtc/StunServers.svelte';
	import IsWindows from '$lib/detection/IsWindows.svelte';
	import ViGEmDownload from './ViGEmDownload.svelte';
	import { _ } from 'svelte-i18n';
	import LanguageSelector from './LanguageSelector.svelte';
	import TurnServers from '$lib/webrtc/TurnServers.svelte';
</script>

<div class="grid gap-40 w-full md:grid-cols-3 py-24">
	<LanguageSelector />
	<StunServers />
	<TurnServers />

	<IsWindows>
		<ViGEmDownload />
	</IsWindows>
</div>
