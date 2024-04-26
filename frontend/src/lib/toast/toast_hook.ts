import { writable } from 'svelte/store';

interface Toast {
	show: boolean;
	message: string;
	type: ToastType;
}

export enum ToastType {
	INFO = 'info',
	SUCCESS = 'success',
	ERROR = 'error'
}

const toastWritable = writable<Toast>({
	show: false,
	message: '',
	type: ToastType.SUCCESS
});

let timer: number | undefined;

export function showToast(message: string, type: Toast['type']) {
	
	setTimeout(() => {
		toastWritable.set({
			show: true,
			message,
			type
		});

		if (timer) clearTimeout(timer);

		timer = window.setTimeout(() => {
			hideToast();
			timer = undefined;
		}, 2000);

	}, 500);

}

export function hideToast() {
	toastWritable.set({
		show: false,
		message: '',
		type: ToastType.SUCCESS
	});
}

export default toastWritable;
