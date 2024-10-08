import { browser } from '$app/environment'
import { init, register, getLocaleFromNavigator } from 'svelte-i18n'

const defaultLocale = 'en'

register('en', () => import('./en.json'))
register('es', () => import('./es.json'))
register('gl', () => import('./gl.json'))
register('ru', () => import('./ru.json'))
register('fr', () => import('./fr.json'))

init({
	fallbackLocale: defaultLocale,
	initialLocale: browser ? new Intl.Locale(getLocaleFromNavigator() ?? defaultLocale).language : defaultLocale,
})