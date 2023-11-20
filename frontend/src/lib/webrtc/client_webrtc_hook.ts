import { showToast, ToastType } from '$lib/hooks/toast';
import { goto } from '$app/navigation';
import { cloneGamepad } from '$lib/gamepad/gamepad_hook';

enum DataChannelLabel {
	Streaming = 'streaming',
	Controller = 'controller'
}

let peerConnection: RTCPeerConnection | undefined;

function initPeerConnection() {
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
	// create a new webrtc peer
	initPeerConnection();

	if (!peerConnection) {
		showToast('Error creating client', ToastType.ERROR);
		return;
	}

	peerConnection.ontrack = (event) => {
		event.streams[0].getTracks().forEach((track) => {
			console.log(track);
		});
	};

	const controllerChannel = peerConnection.createDataChannel(DataChannelLabel.Controller);

	controllerChannel.onopen = () => {
		// Function to send gamepad data through the RTCPeerConnection

		const sendGamepadData = () => {
			const gamepadData = navigator.getGamepads();

			gamepadData.forEach((gamepad) => {
				if (!gamepad) return;

				const serializedData = JSON.stringify(cloneGamepad(gamepad)); // Example: Convert to JSON
				console.log(serializedData);
				controllerChannel.send(serializedData);
			});

			if (gamepadData) {
				// Convert the data to a string or ArrayBuffer as needed
			}
		};

		// Game loop using requestAnimationFrame
		const gameLoop = () => {
			sendGamepadData(); // Send gamepad data

			// Continue the loop
			requestAnimationFrame(gameLoop);
		};

		// Start the game loop
		gameLoop();
	};

	try {
		const offer = await peerConnection.createOffer();

		await peerConnection.setLocalDescription(offer);

		navigator.clipboard.writeText(signalEncode(offer));

		showToast('Client code copied to clipboard', ToastType.SUCCESS);
	} catch (error) {
		showToast('Error creating client', ToastType.ERROR);
	}
}

function ConnectToHostWeb(hostCode: string) {
	const answerResponse = hostCode.split(';');
	const answer: RTCSessionDescription = signalDecode(answerResponse[0]);

	const remoteCandidates: string[] = signalDecode(answerResponse[1]);

	if (!peerConnection) {
		throw new Error('Peer connection not initialized');
	}

	try {
		if (answerResponse.length !== 2) {
			throw new Error('Invalid answer response');
		}

		peerConnection.setRemoteDescription(answer);
		showToast('Connection stablished successfully', ToastType.SUCCESS);
		goto('/mode/client/connection');
	} catch (e) {
		console.error(e);
		showToast('Error connecting to host', ToastType.ERROR);
	}

	for (const candidate of remoteCandidates) {
		peerConnection.addIceCandidate(
			new RTCIceCandidate({
				candidate,
				sdpMid: '',
				sdpMLineIndex: 0
			})
		);
	}
}

function signalEncode<T>(signal: T) {
	// eslint-disable-next-line @typescript-eslint/ban-ts-comment
	//@ts-ignore
	return window.signalEncode(JSON.stringify(signal));
}

function signalDecode<T>(signal: string) {
	// eslint-disable-next-line @typescript-eslint/ban-ts-comment
	//@ts-ignore
	return JSON.parse(window.signalDecode(signal)) as T;
}
