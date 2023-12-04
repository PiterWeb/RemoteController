# ![Gamepad](https://github.com/PiterWeb/RemoteController/blob/e1d7f45a407cf8ba3d4fedf6f7bb99faf7ee3f88/frontend/src/lib/assets/gamepad.svg) Remote Controller 
Remote gamepads for Windows (Tested on Windows 10/11)



### Website (On construction 🚧): https://remote-controller.vercel.app/ 

### Installation 📦

- Setup latest ViGEm Bus : https://github.com/nefarius/ViGEmBus/releases
- Download the latest Remote Controller executable for your platform : https://github.com/PiterWeb/RemoteController/releases/tag/release
- ✅ You are ready to use it

### Use cases ✨

- Play with friends online
- Controll your games from other windows machines with a gamepad
- Create a gaming cloud platform based on windows server (it would require some modifications to interact with the shell instead of the UI)

### Features 🧩

- [x] Simple UI
- [x] P2P "Decentralized"
- [x] Remote Gamepad
- [ ] Remote Streaming

### How it works 👷‍♂️

This desktop APP is based on the WebRTC 🎞 standard and it uses the power of Go to communicate 🗣 with the Windows APIs (XInput & Windows) and the ViGEm DLL

It uses Wails(Go) & Sveltekit(TS) bringing a powerfull connection between Low-Level Logic and UI

### Thanks to the ViGEm project ♥
ViGEm is making this project a reallity. We embed ViGEm Installation Wizard and ViGEm Client DLLS within the executable for Windows