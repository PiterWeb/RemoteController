import { showToast, ToastType } from '$lib/hooks/toast';
import { goto } from '$app/navigation';
// import { gzip, ungzip } from 'pako';

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

async function CreateClientWeb(streamingPipe: (data: unknown) => void) {
	// create a new webrtc peer
	initPeerConnection();

	if (!peerConnection) {
		showToast('Error creating client', ToastType.ERROR);
		return;
	}

	peerConnection.ondatachannel = (event) => {
		const dataChannel = event.channel;

		if (dataChannel.label !== DataChannelLabel.Streaming) return;

		dataChannel.onmessage = (event) => {
			console.log(event.data);
			// pipe the data to the media source
			streamingPipe(event.data);
		};

		dataChannel.onopen = () => {
			console.log('data channel opened');
		};

		dataChannel.onclose = () => {
			console.log('data channel closed');
		};

		dataChannel.onerror = () => {
			showToast('Error during streaming', ToastType.ERROR);
		};
	};

	const controllerChannel = peerConnection.createDataChannel(DataChannelLabel.Controller);

	controllerChannel.onopen = () => {
		const gamepads = navigator.getGamepads();

		const numberOfGamepads = gamepads.length;

		if (numberOfGamepads == 0) return;

		// const gamepadsWorker = new Worker('$lib/workers/gamepads.ts');

		// gamepadsWorker.onmessage = (event) => {
		// 	controllerChannel.send(JSON.stringify(event.data));
		// };
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

		console.log(remoteCandidates);

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
