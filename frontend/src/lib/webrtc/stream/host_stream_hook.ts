import { showToast, ToastType } from '$lib/hooks/toast';
import iceServers from "$lib/webrtc/ice_servers";

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

const MIME_TYPE = 'video/webm;codecs=vp9,opus';

export async function startStreaming() {
	try {
		const mediastream = await navigator.mediaDevices.getDisplayMedia({
			video: true,
			audio: true
		});

		const mediaRecorder = new MediaRecorder(mediastream, {
			mimeType: MIME_TYPE,
		});

		mediaRecorder.ondataavailable = (e) => {
			if (e.data && e.data.size > 0) {
				console.log(e.data);
				sendChunk(e.data);
			}
		};

		mediaRecorder.start(500);
	} catch (e) {
		showToast('Error starting streaming', ToastType.ERROR);
	}
}

function sendChunk(chunk: Blob) {
	console.log('Sending chunk', chunk);
}


export async function CreateHostStream() {
	initStreamingPeerConnection();

	if (!peerConnection) {
		showToast('Error creating host', ToastType.ERROR);
		return;
	}

	peerConnection.onicecandidate = (event) => {
		if (event.candidate) {
			console.log('Sending ice candidate', event.candidate);
		}
	}

	
	const offer = await peerConnection.createOffer();

	await peerConnection.setLocalDescription(offer);

	console.log('Sending offer', offer);
}