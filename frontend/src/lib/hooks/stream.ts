import { showToast, ToastType } from '$lib/hooks/toast';
import { EventsEmit, EventsOn } from '$lib/wailsjs/runtime/runtime';

const MIME_TYPE = 'video/webm;codecs=vp9,opus';

export async function startStreaming() {
	try {
		const mediastream = await navigator.mediaDevices.getDisplayMedia({
			video: true,
			audio: true
		});

		const mediaRecorder = new MediaRecorder(mediastream, {
			mimeType: MIME_TYPE,
			videoBitsPerSecond: 200000
		});

		mediaRecorder.ondataavailable = (e) => {
			if (e.data && e.data.size > 0) {
				console.log(e.data);
				sendChunk(e.data);
			}
		};

		mediaRecorder.start(1000);
	} catch (e) {
		showToast('Error starting streaming', ToastType.ERROR);
	}
}

function sendChunk(chunk: Blob) {
	// Enviar el chunk a través de la conexión
	// console.log(chunk);

	chunk.text().then((text) => {
		EventsEmit('send-streaming', text);
	});
}

export function consumeStreaming() {
	const mediaSource = new MediaSource();
	const queue: ArrayBuffer[] = [];

	let sourceBuffer: SourceBuffer;

	const video = document.getElementById('stream-video') as HTMLVideoElement;

	video.src = URL.createObjectURL(mediaSource);

	mediaSource.addEventListener(
		'sourceopen',
		() => {
			video.play().catch((e) => console.error(e));

			if (!MediaSource.isTypeSupported(MIME_TYPE)) return;

			sourceBuffer = mediaSource.addSourceBuffer(MIME_TYPE);

			sourceBuffer.addEventListener('error', (e) => {
				console.error(e);
			});

			sourceBuffer.addEventListener('abort', (e) => {
				console.error(e);
			});

			sourceBuffer.addEventListener('update', () => {
				if (queue.length > 0 && !sourceBuffer.updating) {
					sourceBuffer.appendBuffer(queue.shift()!);
				}
			});
		},
		false
	);

	receiveChunk((buffer) => {
		if (!sourceBuffer) return;

		if (sourceBuffer.updating || queue.length > 0) {
			queue.push(buffer);
		} else {
			sourceBuffer.appendBuffer(buffer);
		}
	});

	mediaSource.addEventListener('sourceended', () => {
		console.log('Source ended');
	});

	mediaSource.addEventListener('sourceclose', () => {
		console.log('Source closed');
	});

	mediaSource.addEventListener('error', (e) => {
		console.error(e);
	});

	mediaSource.addEventListener('sourceopen', () => {
		console.log('Source open');
	});
}

function receiveChunk(cllbk: (buff: ArrayBuffer) => void) {
	// Recibir el chunk a través de la conexión
		EventsOn('receive-streaming', async (stream: string) => {
			const blob = new Blob([stream], { type: MIME_TYPE });

			console.log(URL.createObjectURL(blob));

			cllbk(await blob.arrayBuffer());
		});
}
