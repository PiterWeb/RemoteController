import { showToast, ToastType } from '$lib/toast/toast_hook';
import { EventsEmit, EventsOff, EventsOn } from '$lib/wailsjs/runtime/runtime';
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
		iceServers: [
			...exportStunServers(),
			...exportTurnServers()
		]
	});
}

// const MIME_TYPE = 'video/webm;codecs=vp9,opus';

export async function startStreaming() {
	try {
		const mediastream = await navigator.mediaDevices.getDisplayMedia({
			video: true,
			audio: true
		});

		const recorder = new MediaRecorder(mediastream, {
			mimeType: 'video/webm;codecs=vp9,opus'
		});

		return recorder;
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
				const mediarecorder = await startStreaming();
				if (!mediarecorder) return;
				mediarecorder.stream.getTracks().forEach((track) => peerConnection?.addTrack(track));
				await peerConnection.setLocalDescription(await peerConnection.createAnswer());
				EventsOff("streaming-signal-client")
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
