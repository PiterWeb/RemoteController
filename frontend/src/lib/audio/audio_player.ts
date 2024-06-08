import { writable, get } from "svelte/store";

const defaultVolume = 0.1;
const audioVolumeStore = writable(defaultVolume);

function playAudio(name: string) {

    const basePath = '/sounds/';

    const audio = new Audio(basePath + name + ".mp3");

    audio.volume = get(audioVolumeStore);

    navigator.userActivation.isActive && audio.play();

}


function volumeChange(volume: number) {
    audioVolumeStore.set(volume);
}

export {playAudio, volumeChange};