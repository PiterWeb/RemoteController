# Docs 📘

## Frontend (UI Logic)

### [Frontend Docs](./FRONTEND.md)

### Stack:

- 💻 Sveltekit  (UI Framework)
- ✨ Tailwind  (CSS Framework)
- 💅 DaisyUI  (Tailwind Framework)
- 🔢 WASM (for Go compatibility) 
- 🔠 Svelte i18n (Translations) 
- 📦 Wails (Golang bindings for desktop)

## Backend (Dekstop APP Logic)

### [Backend Docs](./BACKEND.md)

### Stack:

- 💻 Go
- 📦 Wails (Desktop APP)
- 🌐 Pion/Webrtc
- 🎮 ViGEm (binary & dll for gamepad virtualization)

## General 

Remote Controller uses web technologies like WebRTC and MediaDevices (displayMedia).

WebRTC is totally supported by all main desktop/mobile browsers and is also available in different languages (Go included)

The purpose of WebRTC is to make a P2P connection between Host and Client devices to send Gamepad Input using data channels and also captured Video/Audio with media channels. 

DisplayMedia is for capturing video/audio from desktop/aplications and them stream it through WebRTC media channel.