import onwebsite from '$lib/detection/onwebsite';
import { GetCurrentOS } from '$lib/wailsjs/go/desktop/App';

type OS = 'WINDOWS' | 'LINUX' | 'MACOS';

export async function isWindows() {
	if (onwebsite) return false;
	else return ((await GetCurrentOS()) as OS) === 'WINDOWS';
}

export async function isLinux() {
    if (onwebsite) return false;
    else return ((await GetCurrentOS()) as OS) === 'LINUX';
}

export async function isMacOS() {
    if (onwebsite) return false;
    else return ((await GetCurrentOS()) as OS) === 'MACOS';
}