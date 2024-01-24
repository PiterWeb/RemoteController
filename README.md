# ![Gamepad](./frontend/src/lib/assets/gamepad.svg) Remote Controller 
Remote gamepads without difficulties

> [!Note]
> Website: 
> https://remote-controller.vercel.app/ 

![Example Image from the Desktop APP](./assets/example.jpg)

### Use cases ✨

- Play with friends online
- Controll your games from other machines with a gamepad
- Create a professional gaming cloud platform (it would require modifications to interact with through the shell)

### Installation 📦

- https://remote-controller.vercel.app/download/

### Guides 📘

- [Instalation guide](https://remote-controller.vercel.app/info/guides/installation/)
- [How to use](https://remote-controller.vercel.app/info/guides/how-to-use/)

### Features 🧩

- [x] Portable
- [x] Simple & Modern UI
- [x] P2P "Decentralized" (WebRTC)
- [x] Support for PC/XBOX Gamepads (XInput & DirectInput)
- [ ] Support for PlayStation 3/4/5 Gamepads
- [x] Windows Support
- [ ] Linux Support 
- [x] Remote Streaming
- [x] Browser Client
- [ ] Support for keyboard/mouse

### OS Support 💻

| Windows 	| Linux 	| MacOS 	| Browser (Only Client) 	|
|---------	|-------	|-------	|---------	|
| ✔       	| ❌     	| ❌     	| ✔       	|

### How it works 👷‍♂️

This desktop APP is based on the WebRTC 🎞 standard and it uses the power of Go to communicate 🗣 with the Gamepad emulation libraries.
In Windows uses the ViGEm Bus Driver with the ViGEm Client DLL

For the low level actions uses Go.
On the other hand the UI works with Web technologies (WASM, Sveltekit, Tailwind, DaisyUI & Typescript)

### Contributting 🤝

If you are intested to contribute to this project you can follow this [guide](./CONTRIBUTING.md)

### Thanks to the ViGEm project ♥
ViGEm is making this project a reallity. We embed ViGEm Installation Wizard and ViGEm Client DLLS within the executable for Windows

<a href="https://www.producthunt.com/products/remote-controller/reviews?utm_source=badge-product_review&utm_medium=badge&utm_souce=badge-remote&#0045;controller" target="_blank"><img src="https://api.producthunt.com/widgets/embed-image/v1/product_review.svg?product_id=565186&theme=light" alt="Remote&#0032;Controller - Play&#0032;LOCAL&#0032;co&#0045;op&#0032;games&#0032;ONLINE | Product Hunt" style="width: 250px; height: 54px;" width="250" height="54" /></a>
