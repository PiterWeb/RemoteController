const MIME_TYPE = 'video/webm;codecs=vp9,opus';

let sourceBuffer: SourceBuffer | undefined;
const queue: ArrayBuffer[] = [];

let video : HTMLVideoElement;

let firstChunk = true;

export function consumeStreaming() {
	const mediaSource = new MediaSource();

	video = document.getElementById('stream-video') as HTMLVideoElement;
	video.src = URL.createObjectURL(mediaSource);

	mediaSource.addEventListener(
		'sourceopen',
		() => {

			if (!MediaSource.isTypeSupported(MIME_TYPE)) return;

			sourceBuffer = mediaSource.addSourceBuffer(MIME_TYPE);

			sourceBuffer.addEventListener('error', (e) => {
				console.error(e);
			});

			sourceBuffer.addEventListener('abort', (e) => {
				console.error(e);
			});

			sourceBuffer.addEventListener('update', () => {
				if (queue.length > 0 && !sourceBuffer!.updating) {
					sourceBuffer!.appendBuffer(queue.shift()!);
				}
			});
		},
		false
	);

	mediaSource.addEventListener('sourceended', () => {
		console.log('Source ended');
		queue.length = 0; // clear queue
		sourceBuffer = undefined;
	});

	mediaSource.addEventListener('sourceclose', () => {
		console.log('Source closed');
		queue.length = 0; // clear queue
		sourceBuffer = undefined;
	});

	mediaSource.addEventListener('error', (e) => {
		console.error(e);
	});

	mediaSource.addEventListener('sourceopen', () => {
		console.log('Source open');
	});
}

export async function receiveStreamChunk(bufferStr: string) {
	if (!sourceBuffer) return;

	if (firstChunk) {
		firstChunk = false;
		video.play().catch((e) => console.error(e));
	}

	console.log('Received chunk');

	const blob = new Blob([bufferStr], { type: MIME_TYPE });

	const arrayBuffer = await blob.arrayBuffer();

	console.log(arrayBuffer);

	if (sourceBuffer.updating || queue.length > 0) {
		queue.push(arrayBuffer);
	} else {
		sourceBuffer.appendBuffer(arrayBuffer);
	}
}
