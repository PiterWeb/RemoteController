<script>
	import { CreateHost } from '$lib/webrtc/host_webrtc_hook';
	import { showToast, ToastType } from '$lib/toast/toast_hook';
	import { _ } from 'svelte-i18n'


	let code = '';

	function handleConnectToClient() {
		if (code.length < 1) {
			showToast($_('code-is-empty'), ToastType.ERROR);
			return;
		}

		CreateHost(code);
	}
</script>

<h2 class="text-center text-[clamp(2rem,6vw,4.2rem)] font-black leading-[1.1] xl:text-left">
	<span
		class="[&amp;::selection]:text-base-content text-transparent relative col-start-1 row-start-1 bg-clip-text bg-gradient-to-r from-blue-700 via-blue-800 to-gray-900"
		>{$_('host_card_title')}
	</span>
</h2>

<div
	class="mt-12 card bg-white border border-gray-200 rounded-lg shadow dark:bg-gray-800 dark:border-gray-700"
>
	<div class="card-body">
		<p class="text-gray-400">{$_('get-your-host-code')}</p>
		<div class="card-actions gap-4 justify-end items-center flex-col w-full">
			<div class="join">
				<input
					type="text"
					id="first_name"
					class="max-w-xs bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-l-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
					placeholder={$_('paste-here-code')}
					required
					bind:value={code}
				/>
				<button on:click={handleConnectToClient} class="btn btn-primary">{$_('connect-to-client')}</button>
			</div>
		</div>
	</div>
</div>
