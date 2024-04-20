import { get, writable } from 'svelte/store';
import type { ServersConfig, ICEServer } from '$lib/webrtc/ice';

const defaultTurnConfig: ServersConfig = {};

const turnServersStore = writable<ServersConfig>(
	JSON.parse(localStorage.getItem('turnServers') ?? 'false') || defaultTurnConfig
);

turnServersStore.subscribe((stunServers) =>
	localStorage.setItem('turnServers', JSON.stringify(stunServers))
);

function removeServerFromGroup(group: string, url: string) {
	turnServersStore.update((stunServers) => {
		stunServers[group].urls = stunServers[group].urls.filter((server) => server !== url);
		return stunServers;
	});
}

function modifyGroup(name: string, newName?: string, username?: string, credential?: string) {
	if (newName) {
		turnServersStore.update((stunServers) => {
			stunServers[newName] = stunServers[name];
			stunServers[newName].username = username;
			stunServers[newName].credential = credential;
			delete stunServers[name];
			return stunServers;
		});

		return;
	}

	turnServersStore.update((stunServers) => {
		stunServers[name].username = username;
		stunServers[name].credential = credential;
		return stunServers;
	});
}

function addServerToGroup(group: string, url: string) {
	turnServersStore.update((stunServers) => {
		stunServers[group].urls.push('turn:' + url);
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
	turnServersStore.update((stunServers) => {
		return {
			...stunServers,
			...newServer
		};
	});
}

function deleteServerGroup(name: string) {
	turnServersStore.update((stunServers) => {
		delete stunServers[name];
		return stunServers;
	});
}

function exportTurnServers(): ICEServer[] {
	const servers = get(turnServersStore);
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
	turnServersStore,
	addServerToGroup,
	removeServerFromGroup,
	createServerGroup,
	deleteServerGroup,
	modifyGroup,
	exportTurnServers
};
