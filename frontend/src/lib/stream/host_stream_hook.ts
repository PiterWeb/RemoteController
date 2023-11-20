import { showToast, ToastType } from '$lib/hooks/toast';
import { EventsEmit } from '$lib/wailsjs/runtime/runtime';

const MIME_TYPE = 'video/webm;codecs=vp9,opus';

export async function startStreaming() {
	try {
		const mediastream = await navigator.mediaDevices.getDisplayMedia({
			video: true,
			audio: true
		});

		const mediaRecorder = new MediaRecorder(mediastream, {
			mimeType: MIME_TYPE,
		});

		mediaRecorder.ondataavailable = (e) => {
			if (e.data && e.data.size > 0) {
				console.log(e.data);
				sendChunk(e.data);
			}
		};

		mediaRecorder.start(500);
	} catch (e) {
		showToast('Error starting streaming', ToastType.ERROR);
	}
}

function sendChunk(chunk: Blob) {
	// Enviar el chunk a través de la conexión

	chunk.text().then((text) => {
		EventsEmit('send-streaming', text);
	});
}
