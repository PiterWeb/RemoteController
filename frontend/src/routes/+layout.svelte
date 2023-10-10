<script>
	// @ts-nocheck

	import '../app.css';

	import { page } from '$app/stores';

	import PageTransition from '$lib/PageTransition.svelte';
</script>

<PageTransition key={$page.url} duration={1000}>
	<div class="hero min-h-screen bg-base-200">
		<div class="hero-content flex-col">
			<slot />
		</div>
	</div>
</PageTransition>
