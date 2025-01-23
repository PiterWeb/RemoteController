type keyHandler = (keycode: string) => void

export function handleKeyDown(callback: keyHandler) {
	
	const handler = (event: KeyboardEvent) => {
		const key = event.code.toUpperCase();

		if (specialKeys.has(key)) return callback(key + '_1');

		callback(key);
	}

	document.addEventListener('keydown', handler);

	return handler
}

export function unhandleKeyDown(callback: ReturnType<typeof handleKeyDown>) {
	document.removeEventListener("keydown", callback)
}

export function handleKeyUp(callback: keyHandler) {

	const handler = (event: KeyboardEvent) => {
		const key = event.code.toUpperCase();

		if (!specialKeys.has(key)) return;

		callback(key + '_0');
	}

	document.addEventListener('keyup', handler);

	return handler
}

export function unhandleKeyUp(callback: ReturnType<typeof handleKeyUp>) {
	document.removeEventListener("keyup", callback)
}

const specialKeys = new Set([
	'SHIFTLEFT',
	'SHIFTRIGHT',
	'CONTROLLEFT',
	'CONTROLRIGHT',
	'ALTLEFT',
	'ALTRIGHT'
]);
