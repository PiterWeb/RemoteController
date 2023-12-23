<script>
	import onwebsite from '$lib/hooks/onwebsite';
	import { showToast, ToastType } from '$lib/hooks/toast';

	import { ClosePeerConnection } from '$lib/webrtc/client_webrtc_hook';
	import { CancelConnection } from '$lib/webrtc/host_webrtc_hook';

	function handleToast() {
		showToast('You are now disconnected', ToastType.INFO);
	}

	ClosePeerConnection(handleToast);
	CancelConnection(handleToast);

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
<div class="flex gap-4 mt-4">
	{#if !onwebsite}
		<div class="card md:w-96 md:h-52 bg-base-100 shadow-xl">
			<div class="card-body">
				<h2 class="card-title">Host</h2>
				<p>Share & stream your game to clients</p>
				<a href="/mode/host" class="btn btn-primary">Start Host</a>
			</div>
		</div>
	{/if}
	<div class="card md:w-96 md:h-52 bg-base-100 shadow-xl">
		<div class="card-body">
			<h2 class="card-title">Client</h2>
			<p>Plug your gamepad and play through streaming</p>
			<a href="/mode/client" class="btn btn-primary">Create Client</a>
		</div>
	</div>
</div>
