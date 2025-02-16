import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import { SvelteKitPWA } from '@vite-pwa/sveltekit';

export default defineConfig({
	plugins: [
		sveltekit(),
		SvelteKitPWA({
			srcDir: './src',
			scope: '/',
			base: '/',
			strategies: 'generateSW',
			manifest: {
				short_name: 'RemoteController',
				name: 'Remote Controller',
				start_url: '/',
				scope: '/',
				theme_color: '#4b6bfb',
				background_color: '#ffffff',
				icons: [
					{
						src: './lib/assets/gamepad.svg',
						sizes: '512x512',
						type: 'image/svg+xml'
					}
				]
			},
			workbox: {
				globPatterns: ['client/**/*.{js,css,ico,png,svg,webp,woff,woff2}']
			},
			devOptions: {
				enabled: true,
				type: 'module',
				navigateFallback: '/'
			},
		})
	],
	build: {
		sourcemap: true,
	}
});
