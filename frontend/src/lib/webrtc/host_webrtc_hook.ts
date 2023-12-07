import {
	TryCreateHost as createHostFn,
	TryClosePeerConnection as closeConnectionFn
} from '$lib/wailsjs/go/desktop/App';

import { showToast, ToastType } from '$lib/hooks/toast';
import { goto } from '$app/navigation';

let host: boolean = false;

export async function CreateHost(client: string) {
	try {
		const hostCode = await createHostFn(client);

		if (isError(hostCode)) {
			throw new Error(hostCode);
		}

		navigator.clipboard.writeText(hostCode);
		showToast('Host code copied to clipboard', ToastType.SUCCESS);

		// TODO
		// Listen for connection state changes and handle them (Wails events) to redirect to the correct page

		goto('/mode/host/connection');
		host = true;
	} catch (e) {
		showToast('Error creating host', ToastType.ERROR);
	}
}

function isError(err: string) {
	return err.toUpperCase().includes('ERROR');
}

export function CancelConnection(fn?: () => void) {
	if (fn) fn();
	if (host) closeConnectionFn();
	host = false;
}
