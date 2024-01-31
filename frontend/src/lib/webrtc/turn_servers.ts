import { writable } from 'svelte/store';


const turnServersStore = writable<string[]>(
	JSON.parse(localStorage.getItem('turnServers') ?? '[]')
);

turnServersStore.subscribe((iceServers) =>
	localStorage.setItem('turnServers', JSON.stringify(iceServers))
);

export default turnServersStore;
