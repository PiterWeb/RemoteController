export enum FIXED_RESOLUTIONS {
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

export const DEFAULT_MAX_FRAMERATE = 60
export const DEFAULT_IDEAL_FRAMERATE = 30

