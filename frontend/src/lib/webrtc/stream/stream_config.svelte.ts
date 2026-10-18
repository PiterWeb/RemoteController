import { browser } from '$app/environment'
import log from '$lib/logger/logger'

export enum FIXED_RESOLUTIONS {
	resolutionNative = "Native",
	resolution1080p = "1080",
	resolution720p = "720",
	resolution480p = "480",
	resolution360p = "360"
}

export const RESOLUTIONS: Map<FIXED_RESOLUTIONS,{width: number, height: number}> = new Map()

RESOLUTIONS.set(FIXED_RESOLUTIONS.resolution1080p, {width: 1920, height: 1080})
RESOLUTIONS.set(FIXED_RESOLUTIONS.resolution720p,{width: 1280, height: 720})
RESOLUTIONS.set(FIXED_RESOLUTIONS.resolution480p, {width:854, height: 480})
RESOLUTIONS.set(FIXED_RESOLUTIONS.resolution360p, {width: 640, height:360})

export const DEFAULT_MAX_FRAMERATE = 75
export const DEFAULT_IDEAL_FRAMERATE = 60

export const MIN_FRAMES = 25
export const MAX_FRAMES = 145

const DEFAULT_PREFERED_CODECS = ["video/AV1","video/H264","video/VP9","video/VP8"]

export const preferedCodecsOrdered = $state({value: getStoredPreferedCodecsOrdered()})

export const usePreferedCodecsOrderedStorage = () => {
	$effect(() => {
		localStorage.setItem("codecs-list", JSON.stringify(preferedCodecsOrdered.value))
	})
} 

export function restoreDefaultCodecs() {
    preferedCodecsOrdered.value = DEFAULT_PREFERED_CODECS;
}

function getStoredPreferedCodecsOrdered() {

	if (browser) {
		const stored: string[] = JSON.parse(localStorage.getItem("codecs-list") ?? '[]')

		if (stored && stored.length > 0) {
			return stored
		}

	}

	return DEFAULT_PREFERED_CODECS

}

try {
	getSortedVideoCodecs().forEach(codec => {
		console.log(codec.mimeType);
	})
} catch (e) {
  console.log(`Error getting codecs ${e}`, {err: true});
}

export function getSortedVideoCodecs() {

	const codecs = RTCRtpReceiver.getCapabilities("video")?.codecs;

	if (!codecs) return [];

	log(preferedCodecsOrdered.value)

	return codecs.sort((a, b) => {
	  const indexA = preferedCodecsOrdered.value.indexOf(a.mimeType);
	  const indexB = preferedCodecsOrdered.value.indexOf(b.mimeType);
	  const orderA = indexA >= 0 ? indexA : Number.MAX_VALUE;
	  const orderB = indexB >= 0 ? indexB : Number.MAX_VALUE;
	  return orderA - orderB;
	});
}

export const videoSpeedOptimizationEnabled = $state({value: false})
