import { showToast, ToastType } from '$lib/hooks/toast';
import { EventsEmit, EventsOn } from '$lib/wailsjs/runtime/runtime';
import iceServers from '$lib/webrtc/ice_servers';
import type { SignalingData } from '$lib/webrtc/stream/stream_signal_hook';

let peerConnection: RTCPeerConnection | undefined;

function initStreamingPeerConnection() {
	if (peerConnection) {
		peerConnection.close();
	}

	peerConnection = new RTCPeerConnection({
		iceServers: [
			{
				urls: iceServers
			}
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

		return mediastream;
	} catch (e) {
		showToast('Error starting streaming', ToastType.ERROR);
		return undefined;
	}
}

export function StopStreaming() {
	try {
		if (!peerConnection) return;

		peerConnection.close();
		peerConnection = undefined;

		showToast('Streaming stopped', ToastType.SUCCESS);
	} catch (e) {
		showToast('Error stopping streaming', ToastType.ERROR);
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
		console.log('Message received', data);
		const { type, offer, candidate, role } = JSON.parse(data) as SignalingData;

		if (!peerConnection) return;

		if (role !== 'client') return;

		switch (type) {
			case 'candidate':
				peerConnection.addIceCandidate(candidate);
				break;
			case 'offer':
				if (!offer) return;
				await peerConnection.setRemoteDescription(offer);
				// eslint-disable-next-line no-case-declarations
				const mediastream = await startStreaming();
				if (!mediastream) return;
				mediastream.getTracks().forEach((track) => peerConnection?.addTrack(track));
				await peerConnection.setLocalDescription(await peerConnection.createAnswer());
				break;
		}
	});

	peerConnection.onconnectionstatechange = async () => {
		if (!peerConnection) return;

		if (peerConnection.connectionState === 'connected') {
			showToast('Connected', ToastType.SUCCESS);
			return;
		}
	};
}
