<script lang="ts">
	import onwebsite from '$lib/detection/onwebsite';
	import { _ } from 'svelte-i18n';
	import { goto } from '$app/navigation';
	import StartTutorial from '$lib/tutorial/driver';
	import type { Driver } from 'driver.js';

	const TUTORIAL_DELAY = 750;

	let tutorialDriver: Driver;
	let currentStep = 0;

	function goNextTutorial(duration: number = TUTORIAL_DELAY) {
		setTimeout(() => {
			tutorialDriver.moveNext();
			currentStep++;
		}, duration);
	}

	function goPrevTutorial(duration: number = TUTORIAL_DELAY) {
		setTimeout(() => {
			tutorialDriver.movePrevious();
			currentStep--;
		}, duration);
	}

	let tutorialSteps = [
		{
			element: '#tutorial-config-btn',
			popover: {
				title: $_('tutorial_config_title'),
				description: $_('tutorial_config_description'),
				onNextClick: () => {
					goto('/mode/config');
					goNextTutorial();
				}
			}
		},
		{
			element: '#tutorial-language',
			popover: {
				title: $_('tutorial_language_title'),
				description: $_('tutorial_language_description'),
				onNextClick: () => goNextTutorial,
				onPrevClick: () => {
					goto('/');
					goPrevTutorial();
				}
			}
		},
		{
			element: '#tutorial-stun-card',
			popover: {
				title: $_('tutorial_stun_title'),
				description: $_('tutorial_stun_description'),
				onNextClick: () => {
					goto('/mode/config/advanced/stun');
					goNextTutorial();
				}
			}
		},
		{
			element: '#tutorial-group-server',
			popover: {
				title: $_('tutorial_group_server_title'),
				description: $_('tutorial_group_server_description'),
				onPrevClick: () => {
					goto('/mode/config');
					goPrevTutorial();
				},
				onNextClick: () => {
					goto('/mode/config');
					goNextTutorial();
				}
			}
		}
	];
</script>

<button
	on:click={() => {
		tutorialDriver = StartTutorial(tutorialSteps);
	}}
	class="btn btn-primary text-white"
>
	{$_('tutorial_btn')}
</button>

<h2 class="text-center text-[clamp(2rem,6vw,4.2rem)] font-black leading-[1.1] xl:text-left">
	<span
		class="[&amp;::selection]:text-base-content text-transparent relative col-start-1 row-start-1 bg-clip-text bg-gradient-to-r from-blue-700 via-blue-800 to-gray-900"
		>{$_('main_title_choose')}
	</span>
	{$_('main_title_your')}

	<span
		class="[&amp;::selection]:text-base-content text-transparent relative col-start-1 row-start-1 bg-clip-text bg-gradient-to-r from-blue-700 via-blue-800 to-gray-900"
		>{$_('main_title_role')}</span
	>
</h2>
<div class="flex gap-4 mt-4 md:flex-row flex-col">
	{#if !onwebsite}
		<div
			class="card md:w-96 md:h-52 bg-white border border-gray-200 rounded-lg shadow dark:bg-gray-800 dark:border-gray-700"
		>
			<div class="card-body">
				<h2 class="card-title text-white">{$_('host_card_title')}</h2>
				<p class="text-gray-400">{$_('host_card_description')}</p>
				<a href="/mode/host" class="btn btn-primary text-white">{$_('host_card_cta')}</a>
			</div>
		</div>
	{/if}
	<div
		class="card md:w-96 md:h-52 bg-white border border-gray-200 rounded-lg shadow dark:bg-gray-800 dark:border-gray-700"
	>
		<div class="card-body">
			<h2 class="card-title text-white">{$_('client_card_title')}</h2>
			<p class="text-gray-400">{$_('client_card_description')}</p>
			<a href="/mode/client" class="btn btn-primary text-white">{$_('client_card_cta')}</a>
		</div>
	</div>
</div>
