# ![Gamepad](./frontend/src/lib/assets/gamepad.svg) Remote Controller
### An open source & decentralized alternative to Steam remote play (No need to self host)

> [!Note]
> Website: 
> https://remote-controller.vercel.app/ 

## Use cases ✨

- Play with friends online
- Enjoy videogames from other devices using your gamepad
### The limit is in your imagination ... 

![Example Image from the Desktop APP](./assets/example.jpg)

## Installation 📦

- https://remote-controller.vercel.app/download/

## Guides 📘

- [Instalation guide](https://remote-controller.vercel.app/info/guides/installation/)
- [How to use](https://remote-controller.vercel.app/info/guides/how-to-use/)

## Features 🧩

- [x] Portable
- [x] Simple & Modern UI
- [x] P2P "Decentralized" (WebRTC)
- [x] Remote Streaming
- [x] Browser Client
- [ ] Support for keyboard/mouse
- [x] ViGEmBus Setup (Windows)
- [x] Localization (translation to multiple languages)
- [x] Support for STUN & TURN

### OS Support 💻

| Windows 	| Linux 	| MacOS 	| Browser (Only Client) 	|
|---------	|-------	|-------	|---------	|
| ✔       	| ❌     	| ❌     	| ✔ (Known Issues with Safari)       	|

### Gamepad Support 🎮

| PC Controller (XInput/DirectInput) 	| Xbox Controller (XInput) 	| PlayStation Controler
|---------	|-------	|-------	|
| ✔       	| ✔     	| ❌     	|

## Self Hosting ☁

There is no way to self-host the infrastructure of RemoteController because it has no backend. But instead you can self-host if you want the TURN & STUN servers and then add them to the config.

If you want to self-host a TURN server you can [try Coturn](https://github.com/coturn/coturn)

## How it works 👷‍♂️

This desktop APP is based on the WebRTC 🎞 standard and it uses the power of Go to communicate 🗣 with the Gamepad emulation libraries.
In Windows uses the ViGEm Bus Driver with the ViGEm Client DLL

For the low level actions uses Go.
On the other hand the UI works with Web technologies (WASM, Sveltekit, Tailwind, DaisyUI & Typescript)

## Contributting 🤝

If you are interested to contribute to this project you can follow this [guide](./CONTRIBUTING.md)

## Thanks to the ViGEm project ♥
[ViGEmBus](https://github.com/nefarius/ViGEmBus) is making this project a reallity. We embed ViGEmBus Installation Wizard and ViGEmBus Client DLLS within the executable for Windows

## Did you like the project 👍 ?
You can review us on Product Hunt

<a href="https://www.producthunt.com/products/remote-controller/reviews?utm_source=badge-product_review&utm_medium=badge&utm_souce=badge-remote&#0045;controller" target="_blank"><img src="https://api.producthunt.com/widgets/embed-image/v1/product_review.svg?product_id=565186&theme=light" alt="Remote&#0032;Controller - Play&#0032;LOCAL&#0032;co&#0045;op&#0032;games&#0032;ONLINE | Product Hunt" style="width: 250px; height: 54px;" width="250" height="54" /></a>
