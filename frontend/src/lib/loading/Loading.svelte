<script lang="ts">
	import { playAudio } from '$lib/audio/audio_player';
	import loadingWritable from '$lib/loading/loading_hook';

	let loadingModal: HTMLDialogElement;
	let title: string = '';
	let message: string = '';

	$: if ($loadingWritable.loading && loadingModal) {
		loadingModal.showModal();
		title = $loadingWritable.title ?? title;
		message = $loadingWritable.message ?? message;
		playAudio('open_modal');
	} else if (loadingModal) {
		loadingModal.close();
	}
</script>

<dialog class="modal modal-bottom sm:modal-middle" bind:this={loadingModal}>
	<div class="modal-box flex flex-col items-center gap-6">
		<h4 class="font-bold text-lg">{title}</h4>
        <p class="text-lg">{message}</p>
		<span class="loading loading-spinner loading-lg"></span>
	</div>
</dialog>
