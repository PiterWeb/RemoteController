<script>
	import onwebsite from '$lib/detection/onwebsite';
	
</script>

<h2 class="text-center text-[clamp(2rem,6vw,4.2rem)] font-black leading-[1.1] xl:text-left">
	<span
		class="[&amp;::selection]:text-base-content text-transparent relative col-start-1 row-start-1 bg-clip-text bg-gradient-to-r from-blue-700 via-blue-800 to-gray-900"
		>Choose
	</span>
	your

	<span
		class="[&amp;::selection]:text-base-content text-transparent relative col-start-1 row-start-1 bg-clip-text bg-gradient-to-r from-blue-700 via-blue-800 to-gray-900"
		>role</span
	>
</h2>
<div class="flex gap-4 mt-4 md:flex-row flex-col">
	{#if !onwebsite}
		<div
			class="card md:w-96 md:h-52 bg-white border border-gray-200 rounded-lg shadow dark:bg-gray-800 dark:border-gray-700"
		>
			<div class="card-body">
				<h2 class="card-title text-white">Host</h2>
				<p class="text-gray-400">Share & stream your game to clients</p>
				<a href="/mode/host" class="btn btn-primary text-white">Start Host</a>
			</div>
		</div>
	{/if}
	<div
		class="card md:w-96 md:h-52 bg-white border border-gray-200 rounded-lg shadow dark:bg-gray-800 dark:border-gray-700"
	>
		<div class="card-body">
			<h2 class="card-title text-white">Client</h2>
			<p class="text-gray-400">Plug your gamepad and play through streaming</p>
			<a href="/mode/client" class="btn btn-primary text-white">Create Client</a>
		</div>
	</div>
</div>
