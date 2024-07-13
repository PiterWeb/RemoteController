# ![Gamepad](./frontend/src/lib/assets/gamepad.svg) Remote Controller
### An open source & decentralized alternative to Steam remote play (No need to self host)

⌛ Looking for contributions 👈

> [!Note]
> Website: 
> https://remote-controller.vercel.app/ 

## Use cases ✨

- Play with friends online
- Enjoy videogames from other devices using your gamepad
### The limit is in your imagination ... 

![Example Image from the Desktop APP](./assets/example.jpg)

## Download 📦

- https://remote-controller.vercel.app/download/

## Guides 📘

- [Instalation guide](https://remote-controller.vercel.app/info/guides/installation/)
- [How to use](https://remote-controller.vercel.app/info/guides/how-to-use/)

## Resources 📚

- [FAQ](https://remote-controller.vercel.app/info/resources/faq/)
- [Security](https://remote-controller.vercel.app/info/resources/security/)

## Features 🧩

- [x] Portable
- [x] Simple & Modern UI
- [x] P2P "Decentralized" (WebRTC)
- [x] Remote Streaming
- [x] Browser Client
- [x] Support for keyboard
- [x] ViGEmBus Setup (Windows)
- [x] Localization (translation to multiple languages)
- [x] Support for STUN & TURN

### OS Support 💻

| Windows 	| Linux 	| MacOS 	| Browser (Only Client) 	|
|---------	|-------	|-------	|---------	|
| ✔       	| (Only keyboard)⌛ Looking for contributions for Gamepad     	| ❌     	| ✔ (Known Issues with Safari)       	|

### Gamepad Support 🎮

| PC Controller (XInput/DirectInput) 	| Xbox Controller (XInput) 	| PlayStation Controler
|---------	|-------	|-------	|
| ✔       	| ✔     	| ❌     	|

### Translations 🔠

| English 	| Spanish 	| Galician | Russian |Other languages |
|---------	|-------	|-------	| ------- | ------- |
| 100% ✔     	| 100% ✔      	| 100% ✔      	| 29.5% (Not very good level) |⌛ Looking for contributions

## Self Hosting ☁

There is no way to self-host the infrastructure of RemoteController because it has no backend. But instead you can self-host if you want the TURN & STUN servers and then add them to the config.

- If you want to self-host a TURN/STUN server you can [try Coturn](https://github.com/coturn/coturn). (This is only an example, you can choose other STUN/TURN implementations)

- Also you can host the Web version (but it is only frontend) to make like a network of Remote Controller web clients

## Build

### Prerequisites

You must have Wails CLI, NodeJS, npm and Golang installed.

### How to

First go to the frontend folder and run

  `$ npm install`<br>
  `$ npm run build`

Now run the following command on the root directory of the main project:

  `$ wails build`

finally go to the build/bin folder and your executable will be there.

> [!Note]
> Please note the supported platforms in the table

## How it works 👷‍♂️

This desktop APP is based on the WebRTC 🎞 standard and it uses the power of Go to communicate 🗣 with the Gamepad emulation libraries.
In Windows uses the ViGEm Bus Driver with the ViGEm Client DLL

For the low level actions uses Go.
On the other hand the UI works with Web technologies (WASM, Sveltekit, Tailwind, DaisyUI & Typescript)

You can learn more about [how it works](./docs/README.md) under the hood all the project

## Contributting 🤝

If you are interested to contribute to this project you can follow this [guide](./CONTRIBUTING.md)

## Thanks to the ViGEm project ♥
[ViGEmBus](https://github.com/nefarius/ViGEmBus) is making this project a reallity. We embed ViGEmBus Installation Wizard and ViGEmBus Client DLLS within the executable for Windows

## Did you like the project 👍 ?
You can give a star and review us on Product Hunt

<a href="https://www.producthunt.com/products/remote-controller/reviews?utm_source=badge-product_review&utm_medium=badge&utm_souce=badge-remote&#0045;controller" target="_blank"><img src="https://api.producthunt.com/widgets/embed-image/v1/product_review.svg?product_id=565186&theme=light" alt="Remote&#0032;Controller - Play&#0032;LOCAL&#0032;co&#0045;op&#0032;games&#0032;ONLINE | Product Hunt" style="width: 250px; height: 54px;" width="250" height="54" /></a>

## Star History
Here you can see how fast the community is growing
<br/>
[![Star History Chart](https://api.star-history.com/svg?repos=PiterWeb/RemoteController&type=Timeline)](https://star-history.com/#PiterWeb/RemoteController&Timeline)
