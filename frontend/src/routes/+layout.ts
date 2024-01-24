export const ssr = false;
export const prerender = true;

import { browser } from '$app/environment'
import '$lib/i18n/i18n' // Import to initialize. Important :)
import { locale, waitLocale } from 'svelte-i18n'
import type { LayoutLoad } from './$types'

export const load: LayoutLoad = async () => {
	if (browser) {
		locale.set(window.navigator.language)
	}
	await waitLocale()
}