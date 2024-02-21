# MoodyAPI

API Clients/Agents/Server for Moody's Infrasturcture, ~~and a place where I learn various programming languages.~~

Language Ingredients: `C++`, `CMake`, `Golang`, `PostgreSQL`, `Rust`

| Component      | Build Status                                                   |
| -------------- | -------------------------------------------------------------- |
| API Server     | ![Build](../../actions/workflows/build-server.yml/badge.svg)   |
| Notifier Agent | ![Build](../../actions/workflows/build-notifier.yml/badge.svg) |
| PiScreen       | ![Build](../../actions/workflows/build-piscreen.yml/badge.svg) |

## Functionalities

- Notification Handling: Pushing notifications to all clients.
- Dynamic DNS Server.
- Monitoring Screen on RPi: CPU, Memory, systemd serivces, IP Address.

### Notification Handlers

- [Notifier](Notifier/): The notification daemon and sender (WIP), written in Rust.

## LICENSE

I don't think anyone would need these code, but just in case, they are licensed under `GPLv3`

Major Credits:

- [@bitbank2/ss_oled](https://github.com/bitbank2/ss_oled) for the SSOLED implementation
- [@bitbank2/BitBang_I2C](https://github.com/bitbank2/BitBang_I2C) for the actural I2C implementation
  - Although I have removed most of their code.
- [@miekg/exdns](https://github.com/miekg/exdns) for the great DNS server library
