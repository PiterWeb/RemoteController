<script lang="ts">
	import media_css from "$lib/css/media-video.css?raw"
	import 'player.style/microvideo';
	import { consumingStream } from '$lib/webrtc/stream/stream_signal_hook.svelte';

    /** @type {{children?: import('svelte').Snippet}} */
    let {children} = $props()

	let media: HTMLElement | null = $state(null)

	// Apply custom styles to the media element
	$effect(() => {
		if (!media) return

		const shadowRoot = media.shadowRoot
		const styles = document.createElement('style')
		styles.textContent = media_css

		shadowRoot?.appendChild(styles)
	})

</script>

{@render children?.()}

<media-theme-microvideo bind:this={media} class:hidden={!consumingStream.value}>
<!-- svelte-ignore a11y_media_has_caption -->
<video
	slot="media"
	id="stream-video"
	playsinline
	>
</video>
</media-theme-microvideo>