import { showToast, ToastType } from '$lib/toast/toast_hook';
import { EventsEmit, EventsOn } from '$lib/wailsjs/runtime/runtime';
// import { stunServersStore} from '$lib/webrtc/stun_servers';
import { get } from 'svelte/store';
import type { SignalingData } from '$lib/webrtc/stream/stream_signal_hook';
import { _ } from 'svelte-i18n';
import { exportStunServers } from '../stun_servers';
import { exportTurnServers } from '../turn_servers';
// import turnServers from '$lib/webrtc/turn_servers';

let peerConnection: RTCPeerConnection | undefined;

function initStreamingPeerConnection() {
	if (peerConnection) {
		peerConnection.close();
	}

	peerConnection = new RTCPeerConnection({
		iceServers: [...exportStunServers(), ...exportTurnServers()]
	});
}

export async function startStreaming() {
	try {
		const mediastream = await navigator.mediaDevices.getDisplayMedia({
			video: { 
				frameRate: { min:30, max: 60 },
				noiseSuppression: true, 
				autoGainControl: true,
			},
			audio: true
		});

		return mediastream;
	} catch (e) {
		showToast(get(_)('error-starting-streaming'), ToastType.ERROR);
		return undefined;
	}
}

export function StopStreaming() {
	try {
		if (!peerConnection) return;

		peerConnection.close();
		peerConnection = undefined;

		showToast(get(_)('streaming-stopped'), ToastType.SUCCESS);
	} catch (e) {
		showToast(get(_)('error-stopping-streaming'), ToastType.ERROR);
	}
}

export function CreateHostStream() {
	initStreamingPeerConnection();

	if (!peerConnection) {
		throw new Error('Error creating stream');
	}

	peerConnection.onicecandidate = (event) => {
		if (event.candidate) {
			const data: SignalingData = {
				type: 'candidate',
				candidate: event.candidate.toJSON(),
				role: 'host'
			};
			EventsEmit('streaming-signal-server', JSON.stringify(data));
			return;
		}

		console.log('ICE gathering complete');

		const answer = peerConnection?.localDescription?.toJSON();
		const data: SignalingData = {
			type: 'answer',
			answer,
			role: 'host'
		};
		EventsEmit('streaming-signal-server', JSON.stringify(data));
	};

	EventsOn('streaming-signal-client', async (data: string) => {
		if (!peerConnection) return;

		const { type, offer, candidate, role } = JSON.parse(data) as SignalingData;

		if (role !== 'client') return;

		switch (type) {
			case 'candidate':
				peerConnection.addIceCandidate(candidate);
				break;
			case 'offer':
				if (!offer) return;
				await peerConnection.setRemoteDescription(offer);
				// eslint-disable-next-line no-case-declarations
				const stream = await startStreaming();

				stream?.getTracks().forEach((track) => {
					const sender = peerConnection?.addTrack(track, stream);
					if (!sender) return;
					const params = sender.getParameters();
					if (!params.encodings) {
						params.encodings = [{}];
					}
					params.encodings.forEach((_, i) => {
						params.encodings[i].maxBitrate =  5_000_000
						params.encodings[i].maxFramerate = 60
						// params.encodings[i].scaleResolutionDownBy = 1.25
						params.encodings[i].priority = "high"
					})
					
					sender.setParameters(params);
				});

				await peerConnection.setLocalDescription(await peerConnection.createAnswer());
				break;
		}
	});

	peerConnection.onconnectionstatechange = async () => {
		if (!peerConnection) return;

		if (peerConnection.connectionState === 'connected') {
			showToast(get(_)('connected'), ToastType.SUCCESS);
			return;
		}
	};
}
