import { writable } from 'svelte/store';

const loadingWritable = writable<boolean>(false);

export function toogleLoading() {
	loadingWritable.update((loading) => !loading);
}


export default loadingWritable;