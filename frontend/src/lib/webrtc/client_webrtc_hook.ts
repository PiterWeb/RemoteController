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
	initPeerConnection();

	if (!peerConnection) {
		showToast('Error creating client', ToastType.ERROR);
		return;
	}

	const candidateList: string[] = [];

	peerConnection.onicecandidate = (ev) => {
		if (ev.candidate == null) {
			return;
		}

		const desc = peerConnection?.localDescription;

		if (!desc) return;

		candidateList.push(ev.candidate.candidate);
	};

	peerConnection.ontrack = (event) => {
		event.streams[0].getTracks().forEach((track) => {
			console.log(track);
		});
	};

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

		// Game loop using requestAnimationFrame
		const gameLoop = () => {
			sendGamepadData();

			// Continue the loop
			requestAnimationFrame(gameLoop);
		};

		// Start the game loop
		gameLoop();
	};

	peerConnection.onnegotiationneeded = async () => {
		if (!peerConnection) return;

		try {
			const offer = await peerConnection.createOffer();

			await peerConnection.setLocalDescription(offer);

			navigator.clipboard.writeText(signalEncode(offer) + ';' + signalEncode(candidateList));

			showToast('Client code copied to clipboard', ToastType.SUCCESS);
		} catch (error) {
			showToast('Error creating client', ToastType.ERROR);
		}
	};
}

function ConnectToHostWeb(hostCode: string) {
	const answerResponse = hostCode.split(';');
	const answer: RTCSessionDescription = signalDecode(answerResponse[0]);

	const remoteCandidates: RTCIceCandidateInit[] = signalDecode(answerResponse[1]);

	if (!peerConnection) {
		throw new Error('Peer connection not initialized');
	}

	try {
		if (answerResponse.length !== 2) {
			throw new Error('Invalid answer response');
		}

		peerConnection.setRemoteDescription(answer);

		for (const candidate of remoteCandidates) {
			peerConnection.addIceCandidate(candidate);
		}

		showToast('Connection stablished successfully', ToastType.SUCCESS);
		goto('/mode/client/connection');
	} catch (e) {
		console.error(e);
		showToast('Error connecting to host', ToastType.ERROR);
	}
}

// Function WASM (GOLANG)
function signalEncode<T>(signal: T) {
	// eslint-disable-next-line @typescript-eslint/ban-ts-comment
	//@ts-ignore
	return window.signalEncode(JSON.stringify(signal));
}

// Function WASM (GOLANG)
function signalDecode<T>(signal: string) {
	// eslint-disable-next-line @typescript-eslint/ban-ts-comment
	//@ts-ignore
	return JSON.parse(window.signalDecode(signal)) as T;
}
