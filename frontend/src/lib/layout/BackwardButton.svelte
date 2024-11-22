<script lang="ts">

	import { goto } from "$app/navigation";

	interface Props {
		path?: string;
		action?: () => void;
	}

	let { path = '/', action = () => {} }: Props = $props();

	function handleClick() {
		goto(path);
		action();
	}

</script>

<button aria-label="go back" id="tutorial-back-btn" onclick={handleClick} class="fixed top-28 md:right-28 right-20 p-3 rounded-full bg-white shadow-lg">
	<svg
		xmlns="http://www.w3.org/2000/svg"
		class="h-6 w-6"
		fill="none"
		viewBox="0 0 24 24"
		stroke="currentColor"
	>
		<path
			stroke-linecap="round"
			stroke-linejoin="round"
			stroke-width="2"
			d="M10 19l-7-7m0 0l7-7m-7 7h18"
		/>
	</svg>
</button>
