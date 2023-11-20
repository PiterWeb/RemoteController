<script>
	import BackwardButton from '$lib/components/BackwardButton.svelte';

	import { CancelConnection } from '$lib/webrtc/host_webrtc_hook';

	function handleGoBack() {
		CancelConnection();
	}

</script>

<BackwardButton path="/" action={handleGoBack} />

<slot />
