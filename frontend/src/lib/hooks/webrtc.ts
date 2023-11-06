import {
	ConnectToHost as connectToHostFn,
	CreateClient as createClientFn,
	CreateHost as createHostFn,
	CloseConnection as closeConnectionFn
} from '$lib/wailsjs/go/desktop/App';

import { showToast, ToastType } from '$lib/hooks/toast';
import { goto } from '$app/navigation';


export async function CreateHost(client: string) {
	try {
		const hostCode = await createHostFn(client);

		if (isError(hostCode)) {
			throw new Error(hostCode);
		}

		navigator.clipboard.writeText(hostCode);
		showToast('Host code copied to clipboard', ToastType.SUCCESS);
		goto('/mode/host/connection');
	} catch (e) {
		showToast('Error creating host', ToastType.ERROR);
	}
}

export async function CreateClient() {
	try {
		const clientCode = await createClientFn();

		if (isError(clientCode)) {
			throw new Error(clientCode);
		}

		navigator.clipboard.writeText(clientCode);

		showToast('Client code copied to clipboard', ToastType.SUCCESS);
	} catch (e) {
		console.error(e);
		showToast('Error connecting to host', ToastType.ERROR);
	}
}

export async function ConnectToHost(client: string) {
	try {
		const response = await connectToHostFn(client);

		if (isError(response)) {
			throw new Error(response);
		}

		showToast('Connection stablished successfully', ToastType.SUCCESS);
		goto('/mode/client/connection');
	} catch (e) {
		showToast('Error connecting to client', ToastType.ERROR);
	}
}

function isError(err: string) {
	return err.toUpperCase().includes('ERROR');
}

export function CancelConnection(fn?: () => void) {
	if (fn) fn();
	closeConnectionFn();
}


