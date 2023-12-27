<script lang="ts">
	import IsWindows from '$lib/detection/IsWindows.svelte';

	import { OpenViGEmWizard } from '$lib/wailsjs/go/desktop/App';
</script>


	<section
		class="max-w-sm p-6 bg-white border border-gray-200 rounded-lg shadow dark:bg-gray-800 dark:border-gray-700"
	>
		<h5 class="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
			ViGEmBus Driver
		</h5>

		<p class="mb-3 font-normal text-gray-700 dark:text-gray-400">
			You need to have installed to use the Host mode.
		</p>
		<p class="mb-3 font-normal text-gray-700 dark:text-gray-400">
			If you already have it installed, ignore this, but if you skipped the installation, you can
			install it here.
		</p>
		<button
			on:click={OpenViGEmWizard}
			class="w-full justify-center inline-flex items-center px-3 py-2 text-sm font-medium text-center text-white bg-blue-700 rounded-lg hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
		>
			Install Now
			<svg
				xmlns="http://www.w3.org/2000/svg"
				width="24"
				height="24"
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-width="2"
				stroke-linecap="round"
				stroke-linejoin="round"
				class="lucide lucide-monitor-down ml-2"
				><path d="M12 13V7" /><path d="m15 10-3 3-3-3" /><rect
					width="20"
					height="14"
					x="2"
					y="3"
					rx="2"
				/><path d="M12 17v4" /><path d="M8 21h8" /></svg
			>
		</button>
	</section>

