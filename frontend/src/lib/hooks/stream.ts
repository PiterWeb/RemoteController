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
	console.log(chunk);
	EventsEmit('streaming', await chunk.text());
}

export async function consumeStreaming(video: HTMLVideoElement) {
	try {
		const mediaSource = new MediaSource();

		const buffer_queue: ArrayBuffer[] = []

		video.src = URL.createObjectURL(mediaSource);

		mediaSource.addEventListener('sourceopen', () => {
			video.play()
			const sourceBuffer = mediaSource.addSourceBuffer(MIME_TYPE);

			sourceBuffer.addEventListener("update", () => {
				if (buffer_queue.length > 0 && !sourceBuffer.updating) sourceBuffer.appendBuffer(buffer_queue.shift()!)
			})

			mediaSource.addEventListener('error', function(e) { console.log('error: ' + mediaSource.readyState); });

			receiveChunk(async (chunk) => {
				if (sourceBuffer.updating || buffer_queue.length > 0) buffer_queue.push(await chunk.arrayBuffer())
				else sourceBuffer.appendBuffer(await chunk.arrayBuffer());
			});
		}, false);
	} catch (e) {
		showToast('Error consuming streaming', 'error');
	}
}

async function receiveChunk(cllbk: (chunk: Blob) => void) {
	// Recibir el chunk a través de la conexión

	EventsOn('streaming', (chunkStr: string) => {
		const chunk = new Blob([chunkStr], { type: MIME_TYPE });

		cllbk(chunk);
	});
}
