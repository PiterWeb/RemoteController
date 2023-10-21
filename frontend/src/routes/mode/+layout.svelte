<script>
	import BackwardButton from '$lib/components/BackwardButton.svelte';

	import { CancelConnection } from '$lib/hooks/webrtc';

	function handleGoBack() {
		CancelConnection();
	}

</script>

<BackwardButton path="/" action={handleGoBack} />

<slot />
