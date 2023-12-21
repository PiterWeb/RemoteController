import { ToastType, showToast } from '$lib/hooks/toast';
import iceServers from '$lib/webrtc/ice_servers';

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

async function CreateClientStream(signalingChannel: RTCDataChannel) {
	initStreamingPeerConnection();

	if (!peerConnection) {
		showToast('Error creating client', ToastType.ERROR);
		return;
	}

	peerConnection.onicecandidate = (e) => {
		if (!e.candidate) return;
		signalingChannel.send(JSON.stringify({ type: 'candidate', candidate: e.candidate }));
	};

	peerConnection.ontrack = (e) => {
		console.log('Track received', e);
	};

	const offer = await peerConnection.createOffer();

	await peerConnection.setLocalDescription(offer);

	signalingChannel.send(JSON.stringify({ type: ' offer', offer }));

	signalingChannel.onmessage = async (e) => {
		const { type, answer, candidate } = JSON.parse(e.data);

		if (!peerConnection) {
			return;
		}

		switch (type) {
			case 'answer':
				await peerConnection.setRemoteDescription(answer);
				break;
			case 'candidate':
				peerConnection.addIceCandidate(candidate);
				break;
		}
	};

	peerConnection.onconnectionstatechange = () => {
		console.log('Connection state changed', peerConnection?.connectionState);
	};

}

function ClosePeerConnection() {
	if (!peerConnection) return;
	peerConnection.close();
	peerConnection = undefined;
}

export { CreateClientStream, ClosePeerConnection };
