import { showToast, ToastType } from '$lib/toast/toast_hook';
import { goto } from '$app/navigation';
import { handleGamepad } from '$lib/devices/gamepad/gamepad_hook.svelte';
import {
	handleKeyDown,
	handleKeyUp,
	unhandleKeyDown,
	unhandleKeyUp
} from '$lib/devices/keyboard/keyboard_hook.svelte';
import { toogleLoading } from '$lib/loading/loading_hook';
import { CreateClientStream } from '$lib/webrtc/stream/client_stream_hook';
import { get } from 'svelte/store';
import { CloseStreamClientConnection } from '$lib/webrtc/stream/client_stream_hook';
import { _ } from 'svelte-i18n';
import { exportStunServers } from './stun_servers';
import { exportTurnServers } from './turn_servers';
import { getConsumingStream, setConsumingStream } from './stream/stream_signal_hook.svelte';
import log from '$lib/logger/logger';
import LANMode from './lan_mode.svelte';
import { videoSpeedOptimizationEnabled } from './stream/stream_config.svelte';
import { handleClick, handleMove, unhandleClick, unhandleMove } from '$lib/devices/mouse/mouse_hook.svelte';
import type Modal from '$lib/layout/Modal.svelte';

enum DataChannelLabel {
	StreamingSignal = 'streaming-signal',
	Controller = 'controller',
	Keyboard = 'keyboard',
	Mouse = 'mouse'
}

interface CreateClientWebOptions {
  easyConnect: boolean;
	lostCodeModalElement?: Modal
}

interface createClientCodeOptions {
  clipboard: boolean;
	lostCodeModalElement?: Modal
}

let peerConnection: RTCPeerConnection | undefined;

function initPeerConnection() {
	if (peerConnection) {
		peerConnection.close();
	}

	peerConnection = new RTCPeerConnection({
		iceServers: LANMode.enabled ? [] : [...exportStunServers(), ...exportTurnServers()],
	});
}

function CloseClientConnection(fn?: () => void) {
	if (!peerConnection) return;
	if (fn) fn();
	peerConnection.close();
	peerConnection = undefined;
}

function createClientCode(options: createClientCodeOptions) {
	const { clipboard, lostCodeModalElement } = options;

	return new Promise<string>((resolve, reject) => {
		(async () => {
			let clientCode: string = '';

			if (!peerConnection) return reject('No peerConnection defined');

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

						clientCode =
							signalEncode(peerConnection?.localDescription) + ';' + signalEncode(candidates);

						if (clipboard && navigator && navigator.clipboard && navigator.clipboard.writeText) {
							navigator.clipboard
								.writeText(clientCode)
								.then(() => {
									showToast(get(_)('client-code-copied-to-clipboard'), ToastType.SUCCESS);
								})
								.catch(() => {
                  showToast(get(_)('error-copying-client-code-to-clipboard'), ToastType.ERROR);
                  lostCodeModalElement?.openModal()
								});
						} else if (clipboard) {
							showToast(get(_)('error-copying-client-code-to-clipboard'), ToastType.ERROR);
						}

						return resolve(clientCode);
					}

					candidates.push(ev.candidate.toJSON());
				};
			} catch (error) {
				log(error, {err: true});
				showToast(get(_)('error-creating-client'), ToastType.ERROR);
				return reject(error);
			}
		})();
	});
}

async function CreateClientWeb(options: CreateClientWebOptions) {
	const { easyConnect, lostCodeModalElement } = options;

	initPeerConnection();

	if (!peerConnection) {
		showToast(get(_)('error-creating-client'), ToastType.ERROR);
		return;
	}

	// handleStreamAudio(peerConnection)
	
	peerConnection.onconnectionstatechange = handleConnectionState;

	const controllerChannel = peerConnection.createDataChannel(DataChannelLabel.Controller);
	const streamingSignalChannel = peerConnection.createDataChannel(DataChannelLabel.StreamingSignal);
	const keyboardChannel = peerConnection.createDataChannel(DataChannelLabel.Keyboard);
  const mouseChannel = peerConnection.createDataChannel(DataChannelLabel.Mouse);

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

	let keyDownHandler: ReturnType<typeof handleKeyDown>;
	let keyUpHandler: ReturnType<typeof handleKeyUp>;

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
		unhandleKeyDown(keyDownHandler);
		unhandleKeyUp(keyUpHandler);
	};
	
	let clickHandler: ReturnType<typeof handleClick>;
  let moveHandler: ReturnType<typeof handleMove>;
	
	mouseChannel.onopen = () => {
	
    const sendMouseData = (output: ArrayBuffer) => {
      mouseChannel.send(output)
  	};
    
    clickHandler = handleClick(sendMouseData)
    moveHandler = handleMove(sendMouseData)
	
	}
	
	mouseChannel.onclose = () => {
    unhandleClick(clickHandler)
	  unhandleMove(moveHandler)
	}

	controllerChannel.onopen = () => {
		handleGamepad(controllerChannel);
	};

	let intervalSignalChannel: NodeJS.Timeout

	streamingSignalChannel.onopen = () => {
		let activeStream = false;

		intervalSignalChannel = setInterval(() => {
			if (!getConsumingStream() && activeStream) {
				activeStream = false;
				CloseStreamClientConnection();
			}
			if (getConsumingStream() == activeStream) return;

			activeStream = true;
			CloseStreamClientConnection();

			const videoElement = document.getElementById('stream-video') as HTMLVideoElement;

			if (!videoElement) {
				log('video element not found');
				return;
			}

			setConsumingStream(true);
			CreateClientStream(streamingSignalChannel, videoElement, videoSpeedOptimizationEnabled.value);
		}, 500);
	};

	streamingSignalChannel.onclose = () => {
		CloseStreamClientConnection();
        clearInterval(intervalSignalChannel);
	};

	return await createClientCode({ clipboard: !easyConnect, lostCodeModalElement});
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
		log(e, {err: true});
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
			import('$lib/bindings/App').then((obj) => obj.NotifyCreateClient()).catch();
			break;
		case 'disconnected':
			showToast(get(_)('connection-lost'), ToastType.ERROR);
			CloseClientConnection();
			CloseStreamClientConnection();
			goto('/');
			// Inside try-catch cause in browser will not work
			import('$lib/bindings/App').then((obj) => obj.NotifyCloseClient()).catch();
			break;
		case 'failed':
			showToast(get(_)('connection-failed'), ToastType.ERROR);
			CloseClientConnection();
			CloseStreamClientConnection();
			goto('/');
			// Inside try-catch cause in
			import('$lib/bindings/App').then((obj) => obj.NotifyCloseClient()).catch();
			break;
		case 'closed':
			showToast(get(_)('connection-closed'), ToastType.ERROR);
			CloseClientConnection();
			CloseStreamClientConnection();
			goto('/');
			// Inside try-catch cause in browser will not work
			import('$lib/bindings/App').then((obj) => obj.NotifyCloseClient()).catch();
			break;
	}
}

// function handleStreamAudio(pc: RTCPeerConnection) {
  
//     pc.addTransceiver('audio', {direction: 'recvonly'})
    
//     let inboundStream: MediaStream;
    
//     const audioElement = document.getElementById("stream-audio") as HTMLAudioElement
    
//     pc.addEventListener("track", (ev) => {
//       if (ev.streams && ev.streams[0]) {
// 			audioElement.srcObject = ev.streams[0];
// 			audioElement.play();
//   		} else {
//   			if (!inboundStream) {
//   				inboundStream = new MediaStream();
//   				audioElement.srcObject = inboundStream;
//   				audioElement.play();
//   			}
//   			inboundStream.addTrack(ev.track);
//   		}
//     })
  
// } 

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

export { CreateClientWeb, ConnectToHostWeb, CloseClientConnection };
