<script>
	import { beforeNavigate } from '$app/navigation';
	import BackwardButton from '$lib/layout/BackwardButton.svelte';
	import { showToast, ToastType } from '$lib/toast/toast_hook';
	import { _ } from 'svelte-i18n'


	import { ClosePeerConnection } from '$lib/webrtc/client_webrtc_hook';
	import { CancelConnection } from '$lib/webrtc/host_webrtc_hook';
	/** @type {{children?: import('svelte').Snippet}} */
	let { children } = $props();

	function handleToast() {
		showToast($_('you-are-now-disconnected'), ToastType.INFO);
	}

	function closeConnection() {
		ClosePeerConnection(handleToast);
		CancelConnection(handleToast);
	}

	beforeNavigate((navigator) => {
		const nextPathname = navigator.to?.url.pathname ?? '';

		const actualPathname = navigator.from?.url.pathname ?? '';

		// If the user is navigating to the same page or one level up, we don't want to close the connection
		// but if goes one level down, we want to close the connection
		if (actualPathname.includes('/mode/client') && nextPathname.includes(actualPathname)) return;
		if (actualPathname.includes('/mode/host') && nextPathname.includes(actualPathname)) return;
		if (actualPathname.includes('/mode/config')) return;

		// If the user tryes to leave the page, we will show the browser's dialog to confirm the action
		if (confirm($_('are-you-sure-you-want-to-leave'))) {
			closeConnection();
		} else {
			return navigator.cancel();
		}

		closeConnection();
	});
</script>

<BackwardButton path="/" />

{@render children?.()}
