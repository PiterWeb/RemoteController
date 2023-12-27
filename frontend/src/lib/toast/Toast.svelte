<script>
	import toast from '$lib/toast/toast_hook';
</script>

{#if $toast.show}
	<div class="toast toast-end">
		<div class={`shadow-lg px-6 py-3 bg-${$toast.type} font-bold`}>
			<span>{$toast.message}</span>
		</div>
	</div>
{/if}

<style>

    .bg-info {
        background-color: #209cee;
    }

    .bg-success {
        background-color: #23d160;
    }

    .bg-error {
        background-color: #ff3860;
        color: white;
    }
</style>
