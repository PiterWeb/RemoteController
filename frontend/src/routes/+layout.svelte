<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { _ } from 'svelte-i18n';
	
	import PageTransition from '$lib/layout/PageTransition.svelte';
	import Toast from '$lib/toast/Toast.svelte';

	import GamepadSVG from '$lib/assets/gamepad.svg?raw';
	import Loading from '$lib/loading/Loading.svelte';
	import log from '$lib/logger/logger';
	import Tooltip from '$lib/layout/Tooltip.svelte';
	import PWAInstall from '$lib/pwa/PWAInstall.svelte';
	import { onNavigate } from '$app/navigation';
	import Modal from '$lib/layout/Modal.svelte';
	import { IS_RUNNING_EXTERNAL } from '$lib/detection/onwebsite';

	/** @type {{children?: import('svelte').Snippet}} */
	let { children } = $props();

	onMount(() => {
		if ('serviceWorker' in navigator && import.meta.env.VITE_ON_WEBSITE === "true") {
			addEventListener('load', function () {
				try {
					navigator.serviceWorker.register('../service-worker.js');
				} catch (e) {
					log(`Service Worker registration failed: ${e}`, {err: true});
				}
			});
		}

		if (!IS_RUNNING_EXTERNAL) {
			(async () => {
			    try {
    				const { EventsOn, Events } = await import('$lib/events/events');
    
    				EventsOn(Events.error, (err: string) => {
     					errorModal.openModal();
      
     					errorText = err;
    				});
				} catch {}
			})();
		}
	});

	onNavigate((navigation) => {
		if (!document.startViewTransition) return;

		return new Promise((resolve) => {
			document.startViewTransition(async () => {
				resolve();
				await navigation.complete;
			});
		});
	});

	let errorModal: Modal;
	let errorText = $state<string>()
	
</script>

<svelte:head>
	<title>LibreRemotePlay - Web Client</title>
	<meta name="description" content="LibreRemotePlay Web CLient" />
</svelte:head>

<Modal bind:this={errorModal}>
   	<p class="text-xl capitalize">{$_('init-error')}</p>
	<div class="flex gap-2">
	    {errorText}
		<button class="btn btn-neutral">{$_('lost-code-modal-close')}</button>
	</div>
</Modal>

<nav id="navbar" class="navbar bg-primary text-primary-content">
	<div class="flex-1">
		<h1>
				<a href="/" class="btn btn-ghost normal-case text-xl items-start content-center">
					{@html GamepadSVG}

					<div class="hidden md:block">LibreRemotePlay</div>
				</a>
		</h1>
	</div>
	<div class="flex flex-row gap-4">
	    <PWAInstall />
        <a id="btn-config" aria-label="config" href="/mode/config" class="btn btn-ghost">
			<svg
				id="tutorial-config-btn"
				xmlns="http://www.w3.org/2000/svg"
				width="24"
				height="24"
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-width="2"
				stroke-linecap="round"
				stroke-linejoin="round"
				class="lucide lucide-settings"
				><path
					d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z"
				/><circle cx="12" cy="12" r="3" /></svg
			>
		</a>
	</div>
</nav>

<Tooltip ref="btn-config" placement="left">
    <p>{$_('config_title')}</p>
</Tooltip>

<PageTransition key={page.url.toString()}>
	<div class="hero min-h-[calc(100vh-4rem)] bg-gray-900">
		<div class="hero-content flex-col w-full">
			{@render children?.()}
		</div>
	</div>
</PageTransition>

<Toast />

<Loading />

<style>
    
    #navbar {
         view-transition-name: navbar;
    }
    
</style>
