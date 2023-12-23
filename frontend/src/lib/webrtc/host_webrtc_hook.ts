import {
	TryCreateHost as createHostFn,
	TryClosePeerConnection as closeConnectionFn
} from '$lib/wailsjs/go/desktop/App';

import { EventsOn, EventsOnce } from '$lib/wailsjs/runtime/runtime';

import { showToast, ToastType } from '$lib/hooks/toast';
import { goto } from '$app/navigation';
import { toogleLoading, setLoadingMessage, setLoadingTitle } from '$lib/hooks/loading';
import { StopStreaming } from '$lib/webrtc/stream/host_stream_hook';

let host: boolean = false;

enum ConnectionState {
	Connected = 'CONNECTED',
	Failed = 'FAILED',
	Disconnected = 'DISCONNECTED'
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

			switch (state.toUpperCase()) {
				case ConnectionState.Connected:
					showToast('Connected', ToastType.SUCCESS);
					host = true;
					goto('/mode/host/connection');
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
	if (!host) return;
	closeConnectionFn();
	if (fn) fn();
	host = false;
	StopStreaming();
}

export function ListenForConnectionChanges() {
	const connectionStateCancelEventListener = EventsOn(
		'connection_state',
		(state: ConnectionState) => {
			switch (state.toUpperCase()) {
				case ConnectionState.Connected:
					showToast('Connected', ToastType.SUCCESS);
					host = true;
					goto('/mode/host/connection');
					break;
				case ConnectionState.Failed:
					showToast('Connection failed', ToastType.ERROR);
					goto('/');
					break;
				case ConnectionState.Disconnected:
					showToast('Connection lost', ToastType.ERROR);
					host = false;
					goto('/');
					connectionStateCancelEventListener();
					break;
			}
		}
	);
}
