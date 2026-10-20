import { exportStunServers } from '$lib/webrtc/stun_servers';
import { setConsumingStream, type SignalingData } from '$lib/webrtc/stream/stream_signal_hook.svelte';
import { exportTurnServers } from '$lib/webrtc/turn_servers';
import LANMode from '$lib/webrtc/lan_mode.svelte';
import log from '$lib/logger/logger';

let peerConnection: RTCPeerConnection | undefined;
let inboundStream: MediaStream | null = null;

function initStreamingPeerConnection() {
	if (peerConnection) {
		peerConnection.close();
	}

	peerConnection = new RTCPeerConnection({
		iceServers: LANMode.enabled ? [] : [...exportStunServers(), ...exportTurnServers()]
	});
}

async function CreateClientStream(
	signalingChannel: RTCDataChannel,
	videoElement: HTMLVideoElement,
	videoSpeedOptimizationEnabled: boolean = false
) {
	initStreamingPeerConnection();

	if (!videoElement || !peerConnection) throw new Error('Error creating stream');

	// Trigger host streaming button
	signalingChannel.send(JSON.stringify({}))

	peerConnection.onconnectionstatechange = () => {
		if (!peerConnection) return;

		const connectionTerminatedOptions: RTCPeerConnectionState[] = ["disconnected", "failed", "closed"]

		if (connectionTerminatedOptions.includes(peerConnection.connectionState)) {
			CloseStreamClientConnection()
		}
	};

	peerConnection.onicecandidate = (event) => {
		if (event.candidate) {

			const data: SignalingData = {
				type: 'candidate',
				candidate: event.candidate.toJSON(),
				role: 'client'
			};
			
			signalingChannel.send(JSON.stringify(data));
			return;
		}

		if (!peerConnection?.currentLocalDescription) return
		
		const data: SignalingData = {
			type: 'answer',
			answer: peerConnection?.currentLocalDescription,
			role: 'client'
		};

		signalingChannel.send(JSON.stringify(data))

	};

	let playingStream = false;

	peerConnection.ontrack = (ev) => {

		if (playingStream) return;

		if (videoSpeedOptimizationEnabled) {
		  ev.receiver.jitterBufferTarget = 0;

  		if ("playoutDelayHint" in ev.receiver) {
        ev.receiver.playoutDelayHint = 0;
  		}
		}
		
		if (ev.streams && ev.streams[0]) {
			ev.streams[0].getTracks().forEach(t => {
			  if (videoSpeedOptimizationEnabled && t.kind == "video" && "contentHint" in t) t.contentHint = "motion";
			  t.addEventListener("ended", () => {CloseStreamClientConnection()}, true)
			})
			videoElement.srcObject = ev.streams[0];
			videoElement.play();
			playingStream = true;
		} else {
			if (!inboundStream) {
				inboundStream = new MediaStream();
				inboundStream
				videoElement.srcObject = inboundStream;
				videoElement.play();
				playingStream = true;
			}
			ev.track.addEventListener("ended", () => {CloseStreamClientConnection()}, true)
			inboundStream.addTrack(ev.track);
			inboundStream.getTracks().forEach(t => {
        if (videoSpeedOptimizationEnabled && t.kind == "video" && "contentHint" in t) t.contentHint = "motion";
				t.addEventListener("ended", () => {CloseStreamClientConnection()}, true)
			})
		}
	};

	let offerArrived = false

	signalingChannel.onmessage = async (e) => {
		if (!peerConnection) return;

		console.log(e.data)
		
		try {
  		const { type, offer, candidate, role } = JSON.parse(e.data) as SignalingData;
  
  		if (role !== 'host') return;
  
  		if (type == "candidate") {
  			try {await peerConnection?.addIceCandidate(candidate)} catch {/** */}
  			return
  		}
  
  		if (type !== 'offer') return;
  		if (!offer || offerArrived) return;
  		
  		offerArrived = true;


			await peerConnection?.setRemoteDescription(offer);
            const answer = await peerConnection?.createAnswer();
            await peerConnection?.setLocalDescription(answer);

		} catch (e) {
			// TODO: manage error
			log(e, {err: true})
			return
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
