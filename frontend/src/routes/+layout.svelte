<script>
	// @ts-nocheck

	import '../app.css';

	import { page } from '$app/stores';

	import PageTransition from '$lib/components/PageTransition.svelte';
	import Toast from '$lib/components/Toast.svelte';

	import GamepadSVG from '$lib/assets/gamepad.svg?raw';
	import Loading from '$lib/components/Loading.svelte';
	import onwebsite from '$lib/hooks/onwebsite';
</script>

<svelte:head>
	<title>Remote Controller - Web Client</title>
	<meta name="description" content="Remote Controller Web CLient" />
</svelte:head>

<nav class="navbar bg-primary text-primary-content">
	<div class="flex-1">
		<h1>
			{#if onwebsite}
				<a href="https://remote-controller.vercel.app/" class="btn btn-ghost normal-case text-xl items-start content-center">
					{@html GamepadSVG}

					<div class="hidden md:block">Remote Controller</div>
				</a>
			{:else}
				<a href="/" class="btn btn-ghost normal-case text-xl items-start content-center">
					{@html GamepadSVG}

					<div class="hidden md:block">Remote Controller</div>
				</a>
			{/if}
		</h1>
	</div>
</nav>

<PageTransition key={$page.url} duration={750}>
	<div class="hero min-h-[calc(100vh-4rem)] bg-base-200">
		<div class="hero-content flex-col">
			<slot />
		</div>
	</div>
</PageTransition>

<Toast />

<Loading />
