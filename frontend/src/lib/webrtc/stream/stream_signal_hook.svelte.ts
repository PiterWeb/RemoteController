export interface SignalingData {
	type: 'offer' | 'answer' | 'candidate';
	offer?: RTCSessionDescriptionInit;
	answer?: RTCSessionDescriptionInit;
	candidate?: RTCIceCandidateInit;
	role: 'host' | 'client';
}

export const consumingStream = $state({value:false});

export const mediaStreams: {value: MediaStream[]} = $state({value: []})