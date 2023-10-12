import {
	ConnectToClient as connectToClientFn,
	ConnectToHost as connectToHostFn,
	CreateHost as createHostFn
} from './wailsjs/go/desktop/App';

import { showToast } from '$lib/hooks/toast';
import { goto } from '$app/navigation';

export async function CreateHost() {
	try {
		const hostCode = await createHostFn();

		if (isError(hostCode)) {
			throw new Error(hostCode);
		}

		navigator.clipboard.writeText(hostCode);
		showToast('Host code copied to clipboard', 'success');
	} catch (e) {
		showToast('Error creating host', 'error');
	}
}

export async function ConnectToHost(host: string) {
	try {
		const clientCode = await connectToHostFn(host);

		if (isError(clientCode)) {
			throw new Error(clientCode);
		}

		navigator.clipboard.writeText(clientCode);

		showToast('Client code copied to clipboard', 'success');
	} catch (e) {
		showToast('Error connecting to host', 'error');
	}
}

export async function ConnectToClient(client: string) {
	try {
		const response = await connectToClientFn(client);

		if (isError(response)) {
			throw new Error(response);
		}

		showToast('Connection stablished successfully', 'success');

		goto('/host/connection');
	} catch (e) {
		showToast('Error connecting to client', 'error');
	}
}

function isError(err: string) {
	return err.toUpperCase().includes('ERROR');
}
