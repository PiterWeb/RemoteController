import { writable } from 'svelte/store';

const defaultStunServers = [
	'stun:stun.l.google.com:19302',
	'stun:stun.ipfire.org:3478',
	'stun:stun.l.google.com:19305'
];

const stunServersStore = writable<string[]>(
	JSON.parse(localStorage.getItem('stunServers') ?? 'false') || defaultStunServers
);

stunServersStore.subscribe((iceServers) =>
	localStorage.setItem('stunServers', JSON.stringify(iceServers))
);

export default stunServersStore;
