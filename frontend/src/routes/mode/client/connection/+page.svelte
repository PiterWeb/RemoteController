<script lang="ts">
	let consuming = false;

	function connectToStream() {
		consuming = true;
	}
</script>

<button class="btn btn-primary" disabled={consuming} on:click={connectToStream}>
	Connect to stream</button
>

{#if consuming}
	<video id="stream-video" controls controlslist="nodownload">
		<track kind="captions" />
	</video>
{/if}