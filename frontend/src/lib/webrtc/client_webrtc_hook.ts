import { showToast, ToastType } from '$lib/toast/toast_hook';
import { goto } from '$app/navigation';
import { cloneGamepad } from '$lib/gamepad/gamepad_hook';
import { toogleLoading } from '$lib/loading/loading_hook';
import { CreateClientStream } from '$lib/webrtc/stream/client_stream_hook';
import stunServers from '$lib/webrtc/stun_servers';
import { streamingConsuming } from './stream/stream_signal_hook';
import { get } from 'svelte/store';

enum DataChannelLabel {
	StreamingSignal = 'streaming-signal',
	Controller = 'controller'
}

let peerConnection: RTCPeerConnection | undefined;

function initPeerConnection() {
	if (peerConnection) {
		peerConnection.close();
	}

	peerConnection = new RTCPeerConnection({
		iceServers: [
			{
				urls: get(stunServers)
			}
		]
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
		showToast('Error creating client', ToastType.ERROR);
		return;
	}

	peerConnection.onconnectionstatechange = handleConnectionState;

	const controllerChannel = peerConnection.createDataChannel(DataChannelLabel.Controller);
	const streamingSignalChannel = peerConnection.createDataChannel(DataChannelLabel.StreamingSignal);

	controllerChannel.onopen = () => {
		const sendGamepadData = () => {
			const gamepadData = navigator.getGamepads();

			gamepadData.forEach((gamepad) => {
				if (!gamepad) return;

				const serializedData = JSON.stringify(cloneGamepad(gamepad));
				controllerChannel.send(serializedData);
			});
		};

		const gameLoop = () => {
			sendGamepadData();

			// Continue the loop
			requestAnimationFrame(gameLoop);
		};

		// Start the game loop
		gameLoop();
	};

	streamingSignalChannel.onopen = () => {
		const unlistener = streamingConsuming.subscribe((videoElement) => {
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

				navigator.clipboard.writeText(copiedCode);
				showToast('Client code copied to clipboard', ToastType.SUCCESS);
				return;
			}

			candidates.push(ev.candidate.toJSON());
		};
	} catch (error) {
		console.error(error);
		showToast('Error creating client', ToastType.ERROR);
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
		showToast('Error connecting to host', ToastType.ERROR);
	}
}

function handleConnectionState() {
	if (!peerConnection) return;

	const connectionState = peerConnection.connectionState;

	switch (connectionState) {
		case 'connected':
			showToast('Connection stablished successfully', ToastType.SUCCESS);
			goto('/mode/client/connection');
			break;
		case 'disconnected':
			showToast('Connection lost', ToastType.ERROR);
			ClosePeerConnection();
			goto('/');
			break;
		case 'failed':
			showToast('Connection failed', ToastType.ERROR);
			ClosePeerConnection();
			goto('/');
			break;
		case 'closed':
			showToast('Connection closed', ToastType.ERROR);
			ClosePeerConnection();
			goto('/');
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
