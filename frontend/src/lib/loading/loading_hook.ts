import { writable, get } from 'svelte/store';
import { _ } from 'svelte-i18n';

interface LoadingStore {
	loading: boolean;
	title?: string;
	message?: string;
}

const defaultLoadingStore: LoadingStore = {
	loading: false
};

const loadingWritable = writable<LoadingStore>(defaultLoadingStore);

export function toogleLoading() {
	const currentLoading = get(loadingWritable);

	if (!currentLoading.message && !currentLoading.title) {
		loadingWritable.update((store) => ({
			...store,
			title: get(_)('default-loading-title'),
			message: get(_)('default-loading-message')
		}));
		return;
	}

	loadingWritable.update((store) => {
		const loading = !store.loading;

		if (loading) {
			return { ...store, loading };
		}

		return defaultLoadingStore;
	});
}

export function setLoadingTitle(title: string) {
	loadingWritable.update((store) => ({ ...store, title }));
}

export function setLoadingMessage(message: string) {
	loadingWritable.update((store) => ({ ...store, message }));
}

export default loadingWritable;
