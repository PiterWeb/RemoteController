<script>
	import { beforeNavigate } from '$app/navigation';
	import BackwardButton from '$lib/components/BackwardButton.svelte';
	import { showToast, ToastType } from '$lib/toast/toast_hook';
	import { _ } from 'svelte-i18n'


	import { ClosePeerConnection } from '$lib/webrtc/client_webrtc_hook';
	import { CancelConnection } from '$lib/webrtc/host_webrtc_hook';

	function handleToast() {
		showToast($_('you-are-now-disconnected'), ToastType.INFO);
	}

	function closeConnection() {
		ClosePeerConnection(handleToast);
		CancelConnection(handleToast);
	}

	beforeNavigate((navigator) => {

		const url = navigator.to?.url.toString() ?? '';

		if (url.includes('/mode/client') || url.includes('/mode/host')) return;

		closeConnection()

	})

</script>

<BackwardButton path="/"/>

<slot />
