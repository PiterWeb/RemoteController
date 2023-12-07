<script lang="ts">
	import loadingWritable from '$lib/hooks/loading';

	let loadingModal: HTMLDialogElement;

	$: if ($loadingWritable && loadingModal) {
		loadingModal.showModal();
	} else if (loadingModal) {
		loadingModal.close();
	}
</script>

<dialog class="modal modal-bottom sm:modal-middle" bind:this={loadingModal}>
	<div class="modal-box flex flex-col items-center gap-6">
		<h4 class="font-bold text-lg">Waiting to resolve connection!</h4>
        <p class="text-lg">Make sure to stay focused on this window</p>
		<span class="loading loading-spinner loading-lg"></span>
	</div>
</dialog>
