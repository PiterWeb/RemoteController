<script lang="ts">
	import { locale, locales } from 'svelte-i18n';
	import { _ } from 'svelte-i18n';

	function setLocale(event: Event) {
		const target = event.target as HTMLButtonElement;
		locale.set(target.dataset.value);
	}

	function getLocaleName(locale: string) {

        return new Intl.DisplayNames(locale, { type: 'language' }).of(locale) ?? "";
	}
</script>

<section
	class="
    max-w-96 min-w-full p-4 bg-white border border-gray-200 rounded-lg shadow sm:p-6 md:p-8 dark:bg-gray-800 dark:border-gray-700"
>
	<label for="helper-text" class="text-xl font-medium text-gray-900 dark:text-white"
		>{$_('language-selector-title')}</label
	>

	<ul class="flex flex-col justify-between w-full">
		{#each $locales as locale}
			<li class="w-full my-4">
				<button
					on:click={setLocale}
					data-value={locale}
					class="text-white capitalize px-6 py-2 bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm  dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
					>{getLocaleName(locale)}</button
				>
			</li>
		{/each}
	</ul>
</section>
