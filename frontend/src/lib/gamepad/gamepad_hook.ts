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
