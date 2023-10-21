import { showToast } from './toast';
import { EventsEmit, EventsOn } from '$lib/wailsjs/runtime/runtime';

const MIME_TYPE = 'video/webm;codecs=vp9,opus';

export async function startStreaming(audio: boolean) {
	try {
		const mediastream = await navigator.mediaDevices.getDisplayMedia({
			video: true,
			audio
		});

		const mediaRecorder = new MediaRecorder(mediastream, {
			mimeType: MIME_TYPE
		});

		mediaRecorder.ondataavailable = (e) => {
			if (e.data.size > 0) {
				console.log(e.data);
				sendChunk(e.data);
			}
		};

		mediaRecorder.start(100);
	} catch (e) {
		showToast('Error starting streaming', 'error');
	}
}

async function sendChunk(chunk: Blob) {
	// Enviar el chunk a través de la conexión
	// console.log(chunk);
	EventsEmit('send-streaming', await chunk.text());
}

export function consumeStreaming(video: HTMLVideoElement) {
	const mediaSource = new MediaSource();
	const queue: ArrayBuffer[] = [];

	try {
		video.src = URL.createObjectURL(mediaSource);

		mediaSource.addEventListener(
			'sourceopen',
			() => {
				

				const sourceBuffer = mediaSource.addSourceBuffer(MIME_TYPE);

				sourceBuffer.addEventListener('error', (e) => {
					console.error(e);
				});

				sourceBuffer.addEventListener('update', () => {
					if (queue.length > 0 && !sourceBuffer.updating) {
						sourceBuffer.appendBuffer(queue.shift()!);
					}
				});

				receiveChunk(async (buffer) => {

					if (sourceBuffer.updating) {
						queue.push(buffer);
						return;
					}

					sourceBuffer.appendBuffer(buffer);

					if (video.paused) {
						video.play();
					}
				});
			},
			false
		);

		return mediaSource;
	} catch (e) {
		showToast('Error consuming streaming', 'error');
		return null;
	}
}

function receiveChunk(cllbk: (buff: ArrayBuffer) => void) {
	// Recibir el chunk a través de la conexión

	EventsOn('receive-streaming', async (chunkStr: string) => {
		const chunk = new Blob([chunkStr], { type: MIME_TYPE });

		cllbk(await chunk.arrayBuffer());
	});
}
