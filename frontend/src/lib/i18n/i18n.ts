import { browser } from '$app/environment'
import { init, register } from 'svelte-i18n'

const defaultLocale = 'en'

register('en', () => import('./en.json'))
register('es', () => import('./es.json'))
register('gl', () => import('./gl.json'))
register('ru', () => import('./ru.json'))

init({
	fallbackLocale: defaultLocale,
	initialLocale: browser ? window.navigator.language : defaultLocale,
})