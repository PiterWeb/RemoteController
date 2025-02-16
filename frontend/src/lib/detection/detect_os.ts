import onwebsite from '$lib/detection/onwebsite';
import { GetCurrentOS } from '$lib/wailsjs/go/desktop/App';

type OS = 'WINDOWS' | 'LINUX' | 'MACOS';

export async function isWindows() {
    try {
        if (onwebsite) return false;
        else return ((await GetCurrentOS()) as OS) === 'WINDOWS';
    } catch {
        return false
    }
}

export async function isLinux() {
    try {
        if (onwebsite) return false;
        else return ((await GetCurrentOS()) as OS) === 'LINUX';
    } catch {
        return false
    }
}

export async function isMacOS() {
    try {
        if (onwebsite) return false;
        else return ((await GetCurrentOS()) as OS) === 'MACOS';
    } catch {
        return false
    }
}