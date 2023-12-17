import {
	TryCreateHost as createHostFn,
	TryClosePeerConnection as closeConnectionFn
} from '$lib/wailsjs/go/desktop/App';

import { EventsOnce } from '$lib/wailsjs/runtime/runtime';

import { showToast, ToastType } from '$lib/hooks/toast';
import { goto } from '$app/navigation';
import { toogleLoading, setLoadingMessage, setLoadingTitle } from '$lib/hooks/loading';

let host: boolean = false;

enum ConnectionState {
	Connected = 'CONNECTED',
	Failed = 'FAILED'
}

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

		toogleLoading();
		setLoadingMessage('Waiting for client to connect');
		setLoadingTitle('¡Make sure to pass the code to the client!');

		EventsOnce('connection_state', (state: ConnectionState) => {
			toogleLoading();

			switch (state) {
				case ConnectionState.Connected:
					showToast('Connected', ToastType.SUCCESS);
					host = true;
					goto('/mode/host/connected');
					break;
				case ConnectionState.Failed:
					showToast('Connection failed', ToastType.ERROR);
					goto('/');
					break;
				default:
					showToast('Unknown connection state', ToastType.ERROR);
			}
		});

		
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
