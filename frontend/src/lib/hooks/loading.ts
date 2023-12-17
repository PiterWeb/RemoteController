import { writable } from 'svelte/store';

interface LoadingStore {
	loading: boolean;
	title: string;
	message: string;
}

const defaultLoadingStore: LoadingStore = {
	loading: false,
	title: 'Waiting to resolve connection!',
	message: 'Make sure to stay focused on this window'

}

const loadingWritable = writable<LoadingStore>(defaultLoadingStore);

export function toogleLoading() {
	loadingWritable.update(store => {
		const loading = !store.loading;
		
		if (loading) {
			return {...store, loading};
		}

		return defaultLoadingStore;
	
	});
}

export function setLoadingTitle(title: string) {
	loadingWritable.update(store => ({...store, title}));
}

export function setLoadingMessage(message: string) {
	loadingWritable.update(store => ({...store, message}));
}

export default loadingWritable;