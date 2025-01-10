import { writable } from "svelte/store";

export interface SignalingData {
	type: 'offer' | 'answer' | 'candidate';
	offer?: RTCSessionDescriptionInit;
	answer?: RTCSessionDescriptionInit;
	candidate?: RTCIceCandidateInit;
	role: 'host' | 'client';
}

export const streamingConsumingVideoElement = writable<HTMLVideoElement | undefined>(undefined);

export const consumingStream = $state({value:false});