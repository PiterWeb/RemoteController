export function handleKeyDown(callback: (keycode: string) => void) {
	document.addEventListener('keydown', (event) => {
		const key = event.code.toUpperCase();

		if (specialKeys.has(key)) {
			callback(key + '_1');
			return;
		}

		callback(key);
	});
}

export function handleKeyUp(callback: (keycode: string) => void) {
	document.addEventListener('keyup', (event) => {
		const key = event.code.toUpperCase();

		if (specialKeys.has(key)) {
			callback(key + '_0');
		}
	});
}

const specialKeys = new Set([
	'SHIFTLEFT',
	'SHIFTRIGHT',
	'CONTROLLEFT',
	'CONTROLRIGHT',
	'ALTLEFT',
	'ALTRIGHT'
]);
