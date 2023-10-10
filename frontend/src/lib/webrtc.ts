import {
	ConnectToClient as connectToClientFn,
	ConnectToHost as connectToHostFn,
	CreateHost as createHostFn
} from './wailsjs/go/desktop/App';

export async function CreateHost() {
	try {
		const hostCode = await createHostFn();
		navigator.clipboard.writeText(hostCode);
	} catch (e) {
		console.log(e);
	}
}

export async function ConnectToHost(host: string) {
	try {
		const clientCode = await connectToHostFn(host);

		navigator.clipboard.writeText(clientCode);
	} catch (e) {
		console.log(e);
	}
}

export function ConnectToClient(client: string) {
	return connectToClientFn(client);
}
