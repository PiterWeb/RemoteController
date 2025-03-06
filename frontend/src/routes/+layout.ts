export const ssr = false;
export const prerender = true;

import { browser } from '$app/environment'
import '$lib/i18n/i18n' // Import to initialize. Important :)
import { getLocaleFromNavigator, locale, waitLocale } from 'svelte-i18n'
import type { LayoutLoad } from './$types'
import { getLocaleFromLocalStorage } from '$lib/i18n/i18n';

export const load: LayoutLoad = async () => {
	if (browser) {
		locale.set(getLocaleFromLocalStorage() ?? getLocaleFromNavigator())
	}
	await waitLocale()
}