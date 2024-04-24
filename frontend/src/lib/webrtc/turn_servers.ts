import { get, writable } from 'svelte/store';
import type { ServersConfig, ICEServer } from '$lib/webrtc/ice';

const defaultTurnServers: string[] = [];

const defaultTurnConfig: Readonly<ServersConfig> = {
	default: {
		urls: defaultTurnServers
	}
};

const turnServersStore = writable<ServersConfig>(
	JSON.parse(localStorage.getItem('turnServers') ?? 'false') || defaultTurnConfig
);

turnServersStore.subscribe((turnServers) =>
	localStorage.setItem('turnServers', JSON.stringify(turnServers))
);

function removeServerFromGroup(group: string, url: string) {
	turnServersStore.update((turnServers) => {
		turnServers[group].urls = turnServers[group].urls.filter((server) => server !== url);
		return turnServers;
	});
}

function modifyGroup(name: string, newName?: string, username?: string, credential?: string) {
	if (newName) {
		turnServersStore.update((turnServers) => {
			turnServers[newName] = turnServers[name];

			turnServers[newName].username = username;
			turnServers[newName].credential = credential;
			delete turnServers[name];
			return turnServers;
		});

		return;
	}

	turnServersStore.update((turnServers) => {
		turnServers[name].username = username;
		turnServers[name].credential = credential;
		return turnServers;
	});
}

function addServerToGroup(group: string, url: string) {
	turnServersStore.update((turnServers) => {
		turnServers[group].urls.push('turn:' + url);
		return turnServers;
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
	turnServersStore.update((turnServers) => {
		return {
			...turnServers,
			...newServer
		};
	});
}

function deleteServerGroup(name: string) {
	turnServersStore.update((turnServers) => {
		delete turnServers[name];
		return turnServers;
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
	exportTurnServers,
	defaultTurnConfig
};
