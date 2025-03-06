export type ClonedGamepad = {
	axes: number[];
	buttons: GamepadButton[];
	connected: boolean;
	id: string;
	index: number;
};

type GamepadButton = {
	pressed: boolean;
	value: number;
};


export function cloneGamepad(gamepad: Gamepad): ClonedGamepad {

	return {
		axes: [...gamepad.axes],
		buttons: gamepad.buttons.map((button) => {
			return {
				pressed: button.pressed,
				value: button.value
			};
		}),
		connected: gamepad.connected,
		id: gamepad.id,
		index: gamepad.index
	};
}

export function handleGamepad(controllerChannel: RTCDataChannel) {
	const sendGamepadData = () => {
		const gamepadData = navigator.getGamepads();

		gamepadData.forEach((gamepad) => {
			if (!gamepad) return;

			const serializedData = JSON.stringify(cloneGamepad(gamepad));
			controllerChannel.send(serializedData);
		});
	};

	const gamepadLoop = () => {
		sendGamepadData();

		// Continue the loop
		requestAnimationFrame(gamepadLoop);
	};

	// Start the gamepad loop
	gamepadLoop();
}
