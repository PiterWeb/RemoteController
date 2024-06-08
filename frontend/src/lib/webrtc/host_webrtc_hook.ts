import {
	TryCreateHost as createHostFn,
	TryClosePeerConnection as closeConnectionFn
} from '$lib/wailsjs/go/desktop/App';

import { EventsOn, EventsOnce } from '$lib/wailsjs/runtime/runtime';

import { _ } from 'svelte-i18n'
import { get } from 'svelte/store';
import { showToast, ToastType } from '$lib/toast/toast_hook';
import { goto } from '$app/navigation';
import { toogleLoading, setLoadingMessage, setLoadingTitle } from '$lib/loading/loading_hook';
import { StopStreaming } from '$lib/webrtc/stream/host_stream_hook';
import type { ICEServer } from '$lib/webrtc/ice';
import { exportStunServers } from './stun_servers';
import { exportTurnServers } from './turn_servers';

let host: boolean = false;

enum ConnectionState {
	Connected = 'CONNECTED',
	Failed = 'FAILED',
	Disconnected = 'DISCONNECTED'
}

export async function CreateHost(client: string) {
	try {

		const ICEServers: ICEServer[] = [
			...exportStunServers(),
			...exportTurnServers()
		]

		const hostCode = await createHostFn(ICEServers, client);

		if (isError(hostCode)) {
			throw new Error(hostCode);
		}

		if (navigator && navigator.clipboard && navigator.clipboard.writeText) {
			navigator.clipboard.writeText(hostCode).catch(() => {
				showToast(get(_)('error-copying-host-code-to-clipboard'), ToastType.ERROR);
			});
			showToast(get(_)('host-code-copied-to-clipboard'), ToastType.SUCCESS);
		} else {
			showToast(get(_)('error-copying-host-code-to-clipboard'), ToastType.ERROR);
		}

		// TODO
		// Listen for connection state changes and handle them (Wails events) to redirect to the correct page

		toogleLoading();
		setLoadingMessage(get(_)('waiting-for-client-to-connect'));
		setLoadingTitle(get(_)('make-sure-to-pass-the-code-to-the-client'));

		EventsOnce('connection_state', (state: ConnectionState) => {
			toogleLoading();

			switch (state.toUpperCase()) {
				case ConnectionState.Connected:
					showToast(get(_)('connected'), ToastType.SUCCESS);
					host = true;
					goto('/mode/host/connection');
					break;
				case ConnectionState.Failed:
					showToast(get(_)('connection-failed'), ToastType.ERROR);
					goto('/');
					break;
				default:
					showToast(get(_)('unknown-connection-state'), ToastType.ERROR);
			}
		});
	} catch (e) {
		showToast(get(_)('error-creating-host'), ToastType.ERROR);
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
					showToast(get(_)('connected'), ToastType.SUCCESS);
					host = true;
					goto('/mode/host/connection');
					break;
				case ConnectionState.Failed:
					showToast(get(_)('connection-failed'), ToastType.ERROR);
					goto('/');
					break;
				case ConnectionState.Disconnected:
					showToast(get(_)('connection-lost'), ToastType.ERROR);
					host = false;
					goto('/');
					connectionStateCancelEventListener();
					break;
			}
		}
	);
}
