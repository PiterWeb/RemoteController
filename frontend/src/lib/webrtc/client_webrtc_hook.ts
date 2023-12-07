import { showToast, ToastType } from '$lib/hooks/toast';
import { goto } from '$app/navigation';
import { cloneGamepad } from '$lib/gamepad/gamepad_hook';
import { toogleLoading } from '$lib/hooks/loading';

enum DataChannelLabel {
	Streaming = 'streaming',
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
				urls: 'stun:stun.l.google.com:19302'
			}
		]
	});
}

function ClosePeerConnection() {
	if (!peerConnection) return;
	peerConnection.close();
	peerConnection = undefined;
}

export function ClientWebrtc() {
	return {
		CreateClientWeb,
		ConnectToHostWeb,
		ClosePeerConnection
	};
}

async function CreateClientWeb() {
	initPeerConnection();

	if (!peerConnection) {
		showToast('Error creating client', ToastType.ERROR);
		return;
	}

	peerConnection.onconnectionstatechange = handleConnectionState;

	const controllerChannel = peerConnection.createDataChannel(DataChannelLabel.Controller);

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

	// peerConnection.ontrack = (event) => {
	// 	event.streams[0].getTracks().forEach((track) => {
	// 		console.log(track);
	// 	});
	// };

	try {
		const offer = await peerConnection.createOffer();

		await peerConnection.setLocalDescription(offer);

		// Show spinner while waiting for connection
		toogleLoading();

		peerConnection.onicecandidate = (ev) => {
			if (ev.candidate === null) {
				// Disable spinner
				toogleLoading();
				navigator.clipboard.writeText(signalEncode(offer));
				showToast('Client code copied to clipboard', ToastType.SUCCESS);
			}
		};
	} catch (error) {
		showToast('Error creating client', ToastType.ERROR);
	}
}

async function ConnectToHostWeb(hostCode: string) {
	const answer: RTCSessionDescription = signalDecode(hostCode);

	if (!peerConnection) {
		throw new Error('Peer connection not initialized');
	}

	try {
		await peerConnection.setRemoteDescription(answer);
	} catch (e) {
		console.error(e);
		showToast('Error connecting to host', ToastType.ERROR);
	}
}

function handleConnectionState() {
	if (!peerConnection) return;

	const connectionState = peerConnection.connectionState;

	if (connectionState === 'disconnected') {
		showToast('Connection lost', ToastType.ERROR);
		ClosePeerConnection();
		goto('/');
		return;
	}

	if (connectionState === 'failed') {
		showToast('Connection failed', ToastType.ERROR);
		ClosePeerConnection();
		goto('/');
		return;
	}

	if (connectionState === 'closed') {
		showToast('Connection closed', ToastType.ERROR);
		ClosePeerConnection();
		goto('/');
		return;
	}

	if (connectionState === 'connected') {
		showToast('Connection stablished successfully', ToastType.SUCCESS);
		goto('/mode/client/connection');
		return;
	}

	if (connectionState === 'connecting') {
		showToast('Connecting...', ToastType.INFO);
		return;
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
