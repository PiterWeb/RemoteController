import { showToast, ToastType } from '$lib/toast/toast_hook';
import { EventsEmit, EventsOn } from '$lib/wailsjs/runtime/runtime';
import { get } from 'svelte/store';
import { setStreaming, type SignalingData } from '$lib/webrtc/stream/stream_signal_hook.svelte';
import { _ } from 'svelte-i18n';
import { exportStunServers } from '../stun_servers';
import { exportTurnServers } from '../turn_servers';

let peerConnection: RTCPeerConnection | undefined;

function initStreamingPeerConnection() {
	if (peerConnection) {
		peerConnection.close();
	}

	peerConnection = new RTCPeerConnection({
		iceServers: [...exportStunServers(), ...exportTurnServers()]
	});
}


export enum fixedResolutions {
	resolution1080p = "1080",
	resolution720p = "720",
	resolution480p = "480",
	resolution360p = "360"
}

const resolutions: Map<fixedResolutions,{width: number, height: number}> = new Map()

resolutions.set(fixedResolutions.resolution1080p, {width: 1920, height: 1080})
resolutions.set(fixedResolutions.resolution720p,{width: 1280, height: 720})
resolutions.set(fixedResolutions.resolution480p, {width:854, height: 480})
resolutions.set(fixedResolutions.resolution360p, {width: 640, height:360})

let stream: MediaStream | undefined
let unlistenerStreamingSignal: (() => void) | undefined

async function getDisplayMediaStream(resolution: fixedResolutions = fixedResolutions.resolution720p) {
	try {
		const mediastream = await navigator.mediaDevices.getDisplayMedia({
			video: { 
				frameRate: { ideal:30, max: 60 },
				...(resolutions.get(resolution) ?? {}),
				noiseSuppression: true, 
				autoGainControl: true,
			},
			audio: true,
		});

		return mediastream;
	} catch (e) {
		showToast(get(_)('error-starting-streaming'), ToastType.ERROR);
		return undefined;
	}
}

export function StopStreaming() {
	try {
		unlistenerStreamingSignal?.()
		unlistenerStreamingSignal = undefined
		setStreaming(false)
		stream?.getTracks().forEach(t => t.stop()) 

		if (!peerConnection) return;

		peerConnection.close();
		peerConnection = undefined;

		showToast(get(_)('streaming-stopped'), ToastType.SUCCESS);
	} catch (e) {
		showToast(get(_)('error-stopping-streaming'), ToastType.ERROR);
	}
}

export function CreateHostStream(resolution: fixedResolutions = fixedResolutions.resolution720p) {
	initStreamingPeerConnection();

	if (!peerConnection) {
		throw new Error('Error creating stream');
	}

	peerConnection.onconnectionstatechange = async () => {
		if (!peerConnection) return;

		if (peerConnection.connectionState === 'connected') {
			showToast(get(_)('connected'), ToastType.SUCCESS);
			return;
		}

		const connectionTerminatedOptions: RTCPeerConnectionState[] = ["disconnected", "failed", "closed"]

		if (connectionTerminatedOptions.includes(peerConnection.connectionState)) {
			StopStreaming()
			return
		}
	};

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

	let offerArrived = false

	unlistenerStreamingSignal = EventsOn('streaming-signal-client', async (data: string) => {
		if (!peerConnection) return;

		const { type, offer, candidate, role } = JSON.parse(data) as SignalingData;

		if (role !== 'client') return;

		switch (type) {
			case 'candidate':
				try {peerConnection.addIceCandidate(candidate)} catch {/** */}
				break;
			case 'offer':
				if (!offer || offerArrived) return;
				await peerConnection.setRemoteDescription(offer);
				offerArrived = true
				console.log("displaymedia");
				// eslint-disable-next-line no-case-declarations
				stream = await getDisplayMediaStream(resolution);

				stream?.getTracks().forEach((track) => {
					if (!stream) return
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

				stream?.getTracks().forEach(t => t.addEventListener("ended", () => {StopStreaming()}, true) )

				await peerConnection.setLocalDescription(await peerConnection.createAnswer());
				break;
		}
	});

}
