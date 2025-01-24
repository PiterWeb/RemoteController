import { exportStunServers } from '$lib/webrtc/stun_servers';
import { setConsumingStream, type SignalingData } from '$lib/webrtc/stream/stream_signal_hook.svelte';
import { exportTurnServers } from '$lib/webrtc/turn_servers';

let peerConnection: RTCPeerConnection | undefined;
let inboundStream: MediaStream | null = null;

function initStreamingPeerConnection() {
	if (peerConnection) {
		peerConnection.close();
	}

	peerConnection = new RTCPeerConnection({
		iceServers: [...exportStunServers(), ...exportTurnServers()]
	});
}

async function CreateClientStream(
	signalingChannel: RTCDataChannel,
	videoElement: HTMLVideoElement
) {
	initStreamingPeerConnection();

	if (!videoElement || !peerConnection) throw new Error('Error creating stream');

	peerConnection.onconnectionstatechange = () => {
		if (!peerConnection) return;

		const connectionTerminatedOptions: RTCPeerConnectionState[] = ["disconnected", "failed", "closed"]

		if (connectionTerminatedOptions.includes(peerConnection.connectionState)) {
			CloseStreamClientConnection()
		}
	};

	peerConnection.onicecandidate = (e) => {
		if (!e.candidate) return;

		const data: SignalingData = {
			type: 'candidate',
			candidate: e.candidate.toJSON(),
			role: 'client'
		};

		signalingChannel.send(JSON.stringify(data));
	};

	peerConnection.ontrack = (ev) => {
		if (ev.streams && ev.streams[0]) {
			ev.streams[0].getTracks().forEach(t => t.addEventListener("ended", () => {CloseStreamClientConnection()}, true) )
			videoElement.srcObject = ev.streams[0];
			videoElement.play();
		} else {
			if (!inboundStream) {
				inboundStream = new MediaStream();
				videoElement.srcObject = inboundStream;
				videoElement.play();
			}
			ev.track.addEventListener("ended", () => {CloseStreamClientConnection()}, true)
			inboundStream.addTrack(ev.track);
			inboundStream.getTracks().forEach(t => t.addEventListener("ended", () => {CloseStreamClientConnection()}, true))
		}
	};

	const offer = await peerConnection.createOffer({
		offerToReceiveAudio: true,
		offerToReceiveVideo: true
	});

	await peerConnection.setLocalDescription(offer);

	// Configuración de parámetros del códec
	peerConnection.getSenders().forEach((sender) => {
		const params = sender.getParameters();
		if (!params.encodings) {
			params.encodings = [{}];
		}
		params.encodings[0].maxBitrate = 5_000_000; // Configura el bitrate máximo (en bits por segundo)
		params.encodings[0].maxFramerate = 60; // Configura el frame rate máximo
		// params.encodings[i].scaleResolutionDownBy = 1.25
		params.encodings[0].priority = "high"
		sender.setParameters(params);
	});

	const data: SignalingData = {
		type: 'offer',
		offer: offer,
		role: 'client'
	};

	signalingChannel.send(JSON.stringify(data));

	signalingChannel.onmessage = async (e) => {
		const { type, answer, candidate, role } = JSON.parse(e.data) as SignalingData;

		if (!peerConnection) {
			return;
		}

		if (role !== 'host') {
			return;
		}

		switch (type) {
			case 'answer':
				if (!answer) return;
				await peerConnection.setRemoteDescription(answer);
				break;
			case 'candidate':
				try {peerConnection.addIceCandidate(candidate)} catch {/** */}
				break;
		}
	};


}

function CloseStreamClientConnection() {
	setConsumingStream(false)
	if (!peerConnection) return;
	peerConnection.close();
	peerConnection = undefined;
}

export { CreateClientStream, CloseStreamClientConnection };
