<script>

	import { CreateHost, ConnectToClient } from "$lib/webrtc";

	let code = "";

	function handleCreateHost() {
		CreateHost();
	}

	function handleConnectToClient() {

		if (code.length < 1) {
			console.error("Code is empty");
			return
		};

		ConnectToClient(code);
	}

</script>

<h2 class="text-center text-[clamp(2rem,6vw,4.2rem)] font-black leading-[1.1] xl:text-left">
	<span
		class="[&amp;::selection]:text-base-content text-transparent relative col-start-1 row-start-1 bg-clip-text bg-gradient-to-r from-blue-700 via-blue-800 to-gray-900"
		>Host  </span
	>
</h2>

<div class="mt-12 card bg-base-100 shadow-xl">
	<div class="card-body">
		<div class="card-actions gap-4 justify-end items-center flex-col w-full">

			<div class="divider">First Step</div> 
			<button on:click={handleCreateHost} class="btn btn-primary">Create Host</button>

			<div class="divider">Second Step</div> 

			<div class="join">
				<input
					type="text"
					placeholder="Paste here code"
					class="input input-bordered w-full max-w-xs"
					bind:value={code}
				/>
				<button on:click={handleConnectToClient} class="btn btn-primary">Connect to Client</button>
			</div>
		</div>
	</div>
</div>
