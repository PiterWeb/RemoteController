export type ClonedGamepad = {
	axes: number[];
	buttons: GamepadButton[];
	// connected: boolean;
	// id: string;
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
		// connected: gamepad.connected,
		// id: gamepad.id,
		index: gamepad.index
	};
}

export function handleGamepad(controllerChannel: RTCDataChannel, reduceBandwidth = true) {
	
	let channelEnd = false;
	let skipGamepadIter = true;
	
	controllerChannel.addEventListener("close", () => {
		channelEnd = true;
	})
	
	const sendGamepadData = () => {
		
		skipGamepadIter = reduceBandwidth && !skipGamepadIter;
		if (skipGamepadIter) return

		const gamepadData = navigator.getGamepads();

		gamepadData.forEach((gamepad) => {
			if (!gamepad || !gamepad.connected) return;

      const serializedData = JSON.stringify(cloneGamepad(gamepad));
			
      return controllerChannel.send(serializedData);      
		});
	};

	const gamepadLoop = () => {
		sendGamepadData();
	  if (channelEnd) return

		// Continue the loop
		requestAnimationFrame(gamepadLoop);
	};

	// Start the gamepad loop
	gamepadLoop();
}
