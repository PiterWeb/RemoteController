const gamepads = navigator.getGamepads();

for (const gamepad of gamepads) {
    if (!gamepad) continue;

    console.log(gamepad);
}

