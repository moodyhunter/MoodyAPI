# MoodyAPI

API Clients/Agents/Server for Moody's Infrasturcture, ~~and a place where I learn various programming languages.~~

Language Ingredients: `C++`, `CMake`, `Golang`, `PostgreSQL`, `QML`, `Rust`, `TypeScript`

| Component      | Build Status                                                                                                                |
| -------------- | --------------------------------------------------------------------------------------------------------------------------- |
| API Server     | ![Build](../../actions/workflows/build-server.yml/badge.svg)                                                                |
| Web Panel      | ![Build](../../actions/workflows/build-webpanel.yml/badge.svg) ![Lint](../../actions/workflows/lint-webpanel.yml/badge.svg) |
| Qt Client      | ![Build](../../actions/workflows/build-client.yml/badge.svg)                                                                |
| Camera Agent   | ![Build](../../actions/workflows/build-agent.yml/badge.svg)                                                                 |
| Notifier Agent | ![Build](../../actions/workflows/build-notifier.yml/badge.svg)                                                              |

## Functionalities

- Camera Controlling: Whether to turn on. or turn off a remote camera.
- Notification Handling: Pushing notifications to all clients.
- (WIP) A central management panel for all functionalities.
- (TODO) Wireguard Clients Managing: Managing Wireguard clients and setting IPs.
- (TODO) Dynamic DNS for Wireguard Clients.

## Components

- A common [Server](Server/), implemented in Golang
- A [Web Panel](WebPanel/) managing API clients and Wireguard clients, using Next.js (React)
- Some common [assets](assets/), including `systemd` services and some `sudoers` configurations

### Camera Controllers

- [CameraAgent](CameraAgent/): Camera agent monitoring the camera state and perform start/stop tasks, written in Rust
- [Client](Client/): A Qt-based client toggling app, for both Desktop and Android platforms
- [PiScreen](PiScreen/): A C++ SH1106 OLED screen controller to display camera status messages

### Notification Handlers

- [Notifier](Notifier/): The notification daemon and sender (WIP), written in Rust.

## LICENSE

I don't think anyone would need these code, but just in case, they are licensed under `GPLv3`

Credits:

- [@KDAB/android_openssl](https://github.com/KDAB/android_openssl) for prebuilt OpenSSL libraries for Android platforms
- [@bitbank2/ss_oled](https://github.com/bitbank2/ss_oled) for the SSOLED implementation
- [@bitbank2/BitBang_I2C](https://github.com/bitbank2/BitBang_I2C) for the actural I2C implementation
  - Although I have removed most of their code.
