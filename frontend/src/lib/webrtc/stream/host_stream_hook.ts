import { showToast, ToastType } from '$lib/toast/toast_hook';
import { get } from 'svelte/store';
import { setStreaming, type SignalingData } from '$lib/webrtc/stream/stream_signal_hook.svelte';
import { _ } from 'svelte-i18n';
import { exportStunServers } from '../stun_servers';
import { exportTurnServers } from '../turn_servers';
import { DEFAULT_IDEAL_FRAMERATE, DEFAULT_MAX_FRAMERATE, FIXED_RESOLUTIONS, getSortedVideoCodecs, RESOLUTIONS } from './stream_config.svelte';
import log from '$lib/logger/logger';
import LANMode from '$lib/webrtc/lan_mode.svelte';

let peerConnection: RTCPeerConnection | undefined;

function initStreamingPeerConnection() {
	if (peerConnection) {
		peerConnection.close();
	}

	peerConnection = new RTCPeerConnection({
		iceServers: LANMode.enabled ? [] : [...exportStunServers(), ...exportTurnServers()]
	});
}

let stream: MediaStream | undefined
let unlistenerStreamingSignal: (() => () => void) | undefined

async function getDisplayMediaStream(resolution: FIXED_RESOLUTIONS = FIXED_RESOLUTIONS.resolution720p, idealFrameRate = DEFAULT_IDEAL_FRAMERATE, maxFramerate = DEFAULT_MAX_FRAMERATE) {
	try {
		const mediastream = await navigator.mediaDevices.getDisplayMedia({
			video: { 
				frameRate: { ideal: idealFrameRate, max: maxFramerate },
				...(RESOLUTIONS.get(resolution) ?? {}),
				noiseSuppression: false, 
				autoGainControl: false,
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
		unlistenerStreamingSignal?.()()
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

export async function CreateHostStream(resolution: FIXED_RESOLUTIONS = FIXED_RESOLUTIONS.resolution720p, idealFrameRate = DEFAULT_IDEAL_FRAMERATE, maxFramerate = DEFAULT_MAX_FRAMERATE) {
	initStreamingPeerConnection();

	if (!peerConnection) {
		throw new Error('Error creating stream');
	}

	const videoTransceiver = peerConnection.addTransceiver("video", {
		direction: "sendonly",
	});
  	const audioTransceiver = peerConnection.addTransceiver("audio", {
		direction: "sendonly"
	});

  monitorAndAdaptAudioCodec(audioTransceiver.sender)
	videoTransceiver.setCodecPreferences(getSortedVideoCodecs());
  adaptVideoCodec(videoTransceiver.sender)

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

	peerConnection.onicecandidate = async (event) => {
		if (event.candidate) {
			const data: SignalingData = {
				type: 'candidate',
				candidate: event.candidate.toJSON(),
				role: 'host'
			};
			
			const { EventsEmit, Events } = await import('$lib/events/events');
			EventsEmit(Events.streaming_signal_server, JSON.stringify(data));
			return;
		}

	};

	stream = await getDisplayMediaStream(resolution, idealFrameRate, maxFramerate);

	stream?.getTracks().forEach((track) => {
		track.addEventListener(
			'ended',
			() => {
				StopStreaming();
			},
			true
		)
		if (!stream) return;
		peerConnection?.addTrack(track, stream);
	});

	const offer = await peerConnection.createOffer({
        iceRestart: true,
        offerToReceiveAudio: false,
        offerToReceiveVideo: false
  });
	
	await peerConnection.setLocalDescription(offer);

	const data: SignalingData = {
		type: 'offer',
		offer: offer,
		role: 'host'
	};
	
	const { EventsEmit, Events } = await import('$lib/events/events');
	EventsEmit(Events.streaming_signal_server, JSON.stringify(data));

	async function onSignalArrive(data: string) {
		if (!peerConnection) return;

		const { type, answer, candidate, role } = JSON.parse(data) as SignalingData;

		if (role !== 'client') return;

		if (type === "candidate") {
			try {peerConnection.addIceCandidate(candidate)} catch {/** */}
			return
		}

		if (type !== 'answer') return;
		if (!answer) return;

		try {
			await peerConnection.setRemoteDescription(answer);
		} catch (e) {
			// TODO: manage error
			log(e, {err: true})
			return
		}
		
	}

	(async () => {
		const { EventsOn, Events } = await import('$lib/events/events');
		unlistenerStreamingSignal = EventsOn(Events.streaming_signal_client, (data: string) => onSignalArrive(data));
	})()

}

// Example of monitoring and adapting audio codec parameters
function monitorAndAdaptAudioCodec(audioSender: RTCRtpSender) {
	// Monitor network conditions
	setInterval(async () => {
	  const stats = await audioSender.getStats();
	  let availableBitrate = 0;
  
	  // Calculate packet loss using remote-inbound-rtp stats
	  // (outbound-rtp doesn't know about lost packets)
	  stats.forEach(report => {
  		if (report.type === 'candidate-pair' && report.state === 'succeeded') {
  		  availableBitrate = report.availableOutgoingBitrate;
  		}
	  });
  
	  // Adapt Opus parameters based on conditions
	  const parameters = audioSender.getParameters();
  
	  // Find Opus codec in parameters
	  const opusEncodingIdx = parameters.encodings.findIndex(() => 
			parameters.codecs.find(c => c.mimeType.toLowerCase() === 'audio/opus')
	  );
  
	  if (opusEncodingIdx >= 0) {
  		// Adjust bitrate based on available bandwidth
  		if (availableBitrate > 0) {
  		  // Leave headroom for other traffic
  		  const targetBitrate = Math.min(128000, availableBitrate * 0.7);
  		  parameters.encodings[opusEncodingIdx].maxBitrate = targetBitrate;
  		}
  
  		// Note: We're using standard RTCRtpEncodingParameters rather than
  		// non-standard properties like 'networkPriority' which are
  		// Chromium-only and behind flags
    
  		// Apply the changes
  		audioSender.setParameters(parameters);
	  }
	}, 2000); // Check every 2 seconds
}

async function adaptVideoCodec(videoSender: RTCRtpSender) {
  
  const params = videoSender.getParameters()
  
  const encodings: RTCRtpEncodingParameters[] = params.encodings.map(e => {
    return {
      ...e,
      priority: "high",
    }
  })
  
  try {
    await videoSender.setParameters(
      {
        ...params,
        encodings,
        degradationPreference: "maintain-framerate"
      }
    )
  } catch (e) {
    log(e, {err: true})
  }
  
}