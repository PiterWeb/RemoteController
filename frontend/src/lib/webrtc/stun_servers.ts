import { get, writable } from 'svelte/store';
import type { ServersConfig, ICEServer } from '$lib/webrtc/ice';

const defaultStunServers = [
	'stun:stun.l.google.com:19302',
	'stun:stun.ipfire.org:3478',
	'stun:stun.l.google.com:19305'
];

const defaultStunConfig: Readonly<ServersConfig> = {
	default: {
		urls: defaultStunServers
	}
};

const stunServersStore = writable<ServersConfig>(
	JSON.parse(localStorage.getItem('stunServers') ?? 'false') || defaultStunConfig
);

stunServersStore.subscribe((stunServers) =>
	localStorage.setItem('stunServers', JSON.stringify(stunServers))
);

function removeServerFromGroup(group: string, url: string) {
	stunServersStore.update((stunServers) => {
		stunServers[group].urls = stunServers[group].urls.filter((server) => server !== url);
		return stunServers;
	});
}

function modifyGroup(name: string, newName?: string, username?: string, credential?: string) {
	if (newName) {
		stunServersStore.update((stunServers) => {
			stunServers[newName] = stunServers[name];
			stunServers[newName].username = username;
			stunServers[newName].credential = credential;
			delete stunServers[name];
			return stunServers;
		});

		return;
	}

	stunServersStore.update((stunServers) => {
		stunServers[name].username = username;
		stunServers[name].credential = credential;
		return stunServers;
	});
}

function addServerToGroup(group: string, url: string) {
	stunServersStore.update((stunServers) => {
		stunServers[group].urls.push('stun:' + url);
		return stunServers;
	});
}

function createServerGroup(name: string, username?: string, credential?: string) {
	const newServer = {
		[name]: {
			urls: [],
			username: username,
			credential: credential
		}
	};
	stunServersStore.update((stunServers) => {
		return {
			...stunServers,
			...newServer
		};
	});
}

function deleteServerGroup(name: string) {
	stunServersStore.update((stunServers) => {
		delete stunServers[name];
		return stunServers;
	});
}

function exportStunServers(): ICEServer[] {
	const servers = get(stunServersStore);
	const serversArray = Object.keys(servers).map((key) => {
		return {
			urls: servers[key].urls,
			...(servers[key].username && { username: servers[key].username }),
			...(servers[key].credential && { credential: servers[key].credential })
		};
	});

	return serversArray;
}

export {
	stunServersStore,
	addServerToGroup,
	removeServerFromGroup,
	createServerGroup,
	deleteServerGroup,
	modifyGroup,
	exportStunServers,
	defaultStunConfig
};
