import { showToast, ToastType } from '$lib/toast/toast_hook';
import { goto } from '$app/navigation';
import { cloneGamepad } from '$lib/gamepad/gamepad_hook';
import { handleKeyDown, handleKeyUp, unhandleKeyDown, unhandleKeyUp } from '$lib/keyboard/keyboard_hook';
import { toogleLoading } from '$lib/loading/loading_hook';
import { CreateClientStream } from '$lib/webrtc/stream/client_stream_hook';
import { streamingConsumingVideoElement } from './stream/stream_signal_hook';
import { get } from 'svelte/store';
import { CloseStreamPeerConnection } from '$lib/webrtc/stream/client_stream_hook';
import { _ } from 'svelte-i18n';
import { exportStunServers } from './stun_servers';
import { exportTurnServers } from './turn_servers';

enum DataChannelLabel {
	StreamingSignal = 'streaming-signal',
	Controller = 'controller',
	Keyboard = 'keyboard',
}

let peerConnection: RTCPeerConnection | undefined;

function initPeerConnection() {
	if (peerConnection) {
		peerConnection.close();
	}

	peerConnection = new RTCPeerConnection({
		iceServers: [...exportStunServers(), ...exportTurnServers()]
	});
}

function ClosePeerConnection(fn?: () => void) {
	if (!peerConnection) return;
	if (fn) fn();
	peerConnection.close();
	peerConnection = undefined;
}

async function CreateClientWeb() {
	initPeerConnection();

	if (!peerConnection) {
		showToast(get(_)('error-creating-client'), ToastType.ERROR);
		return;
	}

	peerConnection.onconnectionstatechange = handleConnectionState;

	const controllerChannel = peerConnection.createDataChannel(DataChannelLabel.Controller);
	const streamingSignalChannel = peerConnection.createDataChannel(DataChannelLabel.StreamingSignal);
	const keyboardChannel = peerConnection.createDataChannel(DataChannelLabel.Keyboard);

	peerConnection.ondatachannel = (ev) => {
		const channel = ev.channel;

		const label = channel.label;

		channel.onopen = () => {
			console.log('Channel open', label);
		};

		channel.onmessage = (ev) => {
			console.log('Message received', ev.data);
		};
	};

	let keyDownHandler: ReturnType<typeof handleKeyDown>
	let keyUpHandler: ReturnType<typeof handleKeyUp>

	keyboardChannel.onopen = () => {
		const sendKeyboardData = (keycode: string) => {
			console.log('Sending keycode', keycode);
			keyboardChannel.send(keycode);
		};

		// On keydown and keyup events, send the keycode to the host
		keyDownHandler = handleKeyDown(sendKeyboardData);
		keyUpHandler = handleKeyUp(sendKeyboardData);
	};

	keyboardChannel.onclose = () => {
		unhandleKeyDown(keyDownHandler)
		unhandleKeyUp(keyUpHandler)
	}

	controllerChannel.onopen = () => {
		const sendGamepadData = () => {
			const gamepadData = navigator.getGamepads();

			gamepadData.forEach((gamepad) => {
				if (!gamepad) return;

				const serializedData = JSON.stringify(cloneGamepad(gamepad));
				controllerChannel.send(serializedData);
			});
		};

		const gamepadLoop = () => {
			sendGamepadData();

			// Continue the loop
			requestAnimationFrame(gamepadLoop);
		};

		// Start the gamepad loop
		gamepadLoop();
	};

	streamingSignalChannel.onopen = () => {
		const unlistener = streamingConsumingVideoElement.subscribe((videoElement) => {
			if (!videoElement) return;
			CreateClientStream(streamingSignalChannel, videoElement);
			unlistener();
		});
	};

	let copiedCode: string = '';

	try {
		const offer = await peerConnection.createOffer();

		await peerConnection.setLocalDescription(offer);

		// Show spinner while waiting for connection
		toogleLoading();

		const candidates: RTCIceCandidateInit[] = [];

		peerConnection.onicecandidate = (ev) => {
			if (ev.candidate === null) {
				// Disable spinner
				toogleLoading();

				copiedCode =
					signalEncode(peerConnection?.localDescription) + ';' + signalEncode(candidates);

				if (navigator && navigator.clipboard && navigator.clipboard.writeText) {
					navigator.clipboard.writeText(copiedCode).catch(() => {
						showToast(get(_)('error-copying-client-code-to-clipboard'), ToastType.ERROR);
					});
					showToast(get(_)('client-code-copied-to-clipboard'), ToastType.SUCCESS);
				} else {
					showToast(get(_)('error-copying-client-code-to-clipboard'), ToastType.ERROR);
				}
				return;
			}

			candidates.push(ev.candidate.toJSON());
		};
	} catch (error) {
		console.error(error);
		showToast(get(_)('error-creating-client'), ToastType.ERROR);
	}

	return copiedCode;
}

async function ConnectToHostWeb(hostAndCandidatesCode: string) {
	try {
		const [hostCode, candidatesCode] = hostAndCandidatesCode.split(';');

		const answer: RTCSessionDescription = signalDecode(hostCode);

		const candidates: RTCIceCandidateInit[] = signalDecode(candidatesCode);

		if (!peerConnection) {
			throw new Error('Peer connection not initialized');
		}

		await peerConnection.setRemoteDescription(answer);

		candidates.forEach(async (candidate) => {
			if (!peerConnection) return;
			await peerConnection.addIceCandidate(candidate);
		});
	} catch (e) {
		console.error(e);
		showToast(get(_)('error-connecting-to-host'), ToastType.ERROR);
	}
}

function handleConnectionState() {
	if (!peerConnection) return;

	const connectionState = peerConnection.connectionState;

	switch (connectionState) {
		case 'connected':
			showToast(get(_)('connection-established-successfully'), ToastType.SUCCESS);
			goto('/mode/client/connection');
			// Inside try-catch cause in browser will not work
			import('$lib/wailsjs/go/desktop/App').then(obj => obj.NotifyCreateClient()).catch();
			break;
		case 'disconnected':
			showToast(get(_)('connection-lost'), ToastType.ERROR);
			ClosePeerConnection();
			CloseStreamPeerConnection();
			goto('/');
			// Inside try-catch cause in browser will not work
			import('$lib/wailsjs/go/desktop/App').then(obj => obj.NotifyCloseClient).catch();
			break;
		case 'failed':
			showToast(get(_)('connection-failed'), ToastType.ERROR);
			ClosePeerConnection();
			CloseStreamPeerConnection();
			goto('/');
			// Inside try-catch cause in browser will not work
			import('$lib/wailsjs/go/desktop/App').then(obj => obj.NotifyCloseClient).catch();
			break;
		case 'closed':
			showToast(get(_)('connection-closed'), ToastType.ERROR);
			ClosePeerConnection();
			CloseStreamPeerConnection();
			goto('/');
			// Inside try-catch cause in browser will not work
			import('$lib/wailsjs/go/desktop/App').then(obj => obj.NotifyCloseClient).catch();
			break;
	}
}

// Function WASM (GOLANG)
function signalEncode<T>(signal: T): string {
	// eslint-disable-next-line @typescript-eslint/ban-ts-comment
	//@ts-ignore
	return window.signalEncode(JSON.stringify(signal));
}

// Function WASM (GOLANG)
function signalDecode<T>(signal: string): T {
	// eslint-disable-next-line @typescript-eslint/ban-ts-comment
	//@ts-ignore
	return JSON.parse(window.signalDecode(signal)) as T;
}

export { CreateClientWeb, ConnectToHostWeb, ClosePeerConnection };
