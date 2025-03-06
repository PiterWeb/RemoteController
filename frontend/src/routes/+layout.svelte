<script>
	import '../app.css';

	import { page } from '$app/stores';

	import PageTransition from '$lib/layout/PageTransition.svelte';
	import Toast from '$lib/toast/Toast.svelte';

	import GamepadSVG from '$lib/assets/gamepad.svg?raw';
	import Loading from '$lib/loading/Loading.svelte';

	/** @type {{children?: import('svelte').Snippet}} */
	let { children } = $props();

</script>

<svelte:head>
	<title>Remote Controller - Web Client</title>
	<meta name="description" content="Remote Controller Web CLient" />
</svelte:head>

<nav class="navbar bg-primary text-primary-content">
	<div class="flex-1">
		<h1>
				<a href="/" class="btn btn-ghost normal-case text-xl items-start content-center">
					{@html GamepadSVG}

					<div class="hidden md:block">Remote Controller</div>
				</a>
		</h1>
	</div>
	<div class="flex-none">
		<a aria-label="config" href="/mode/config" class="btn btn-ghost">
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

<PageTransition key={$page.url.toString()} duration={750}>
	<div class="hero min-h-[calc(100vh-4rem)] bg-base-200">
		<div class="hero-content flex-col">
			{@render children?.()}
		</div>
	</div>
</PageTransition>

<Toast />

<Loading />
