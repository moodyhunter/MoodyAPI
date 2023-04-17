mod models;
use models::{
    common::Auth,
    dns::{
        CreateDnsRecordRequest, DeleteDnsRecordRequest, DnsRecord, ListDnsRecordsRequest,
        UpdateDnsRecordRequest,
    },
    light::{light_state::Mode, LightColor, LightState, SetLightRequest},
    moody_api::{
        moody_api_service_client::MoodyApiServiceClient, CameraState, UpdateCameraStateRequest,
    },
    notifications::{
        CreateChannelRequest, DeleteChannelRequest, ListChannelRequest, Notification,
        NotificationChannel, SendRequest,
    },
};

use ini::Ini;
use notify_rust::{Hint, Notification as Notify};
use platform_dirs::AppDirs;
use std::{env, future::Future, pin::Pin, process::exit, time::Duration};
use tokio::time::sleep;
use tonic::{transport::Channel, Request, Status};

use crate::models::notifications::SubscribeRequest;

type CommandOp = fn(
    cha: Channel,
    uuid: String,
    args: Vec<String>,
) -> Pin<Box<dyn Future<Output = Result<(), Status>>>>;

struct SubCommand {
    subcommand: &'static str,
    description: &'static str,
    nargs: usize,
    usage: &'static str,
    op: CommandOp,
}

struct Command {
    command: &'static str,
    short: &'static str,
    description: &'static str,
    subcommands: &'static [SubCommand],
}

static COMMANDS: &[Command] = &[
    Command {
        command: "notification",
        short: "n",
        description: "Send to, or subscribe to a notification channel.",
        subcommands: &[
            SubCommand {
                subcommand: "send",
                description: "Send a notification to a channel.",
                nargs: 3,
                usage: "<channel> <title> <content>",
                op: |a, b, c| Box::pin(notification_send(a, b, c)),
            },
            SubCommand {
                subcommand: "subscribe",
                description: "Subscribe to a notification channel.",
                nargs: 1,
                usage: "<channel>",
                op: |a, b, c| Box::pin(notification_subscribe(a, b, c)),
            },
        ],
    },
    Command {
        command: "camera",
        short: "cam",
        description: "Camera operations.",
        subcommands: &[
            SubCommand {
                subcommand: "start",
                description: "Start a camera stream.",
                nargs: 0,
                usage: "",
                op: |a, b, c| Box::pin(camera_start(a, b, c)),
            },
            SubCommand {
                subcommand: "stop",
                description: "Stop a camera stream.",
                nargs: 0,
                usage: "",
                op: |a, b, c| Box::pin(camera_stop(a, b, c)),
            },
        ],
    },
    Command {
        command: "channel",
        short: "chan",
        description: "Create, or delete a notification channel.",
        subcommands: &[
            SubCommand {
                subcommand: "create",
                description: "Create a notification channel.",
                nargs: 1,
                usage: "<channel>",
                op: |a, b, c| Box::pin(channel_create(a, b, c)),
            },
            SubCommand {
                subcommand: "delete",
                description: "Delete a notification channel.",
                nargs: 1,
                usage: "<channel>",
                op: |a, b, c| Box::pin(channel_delete(a, b, c)),
            },
            SubCommand {
                subcommand: "list",
                description: "List all notification channels.",
                nargs: 0,
                usage: "",
                op: |a, b, c| Box::pin(channel_list(a, b, c)),
            },
        ],
    },
    Command {
        command: "dns",
        short: "dns",
        description: "DNS operations.",
        subcommands: &[
            SubCommand {
                subcommand: "update",
                description: "Update the IP address of a DNS record with given domain and type.",
                nargs: 3,
                usage: "<domain> <type> <ip/host>",
                op: |a, b, c| Box::pin(dns_update(a, b, c)),
            },
            SubCommand {
                subcommand: "delete",
                description: "Delete a DNS record with given domain and type.",
                nargs: 2,
                usage: "<domain> <type>",
                op: |a, b, c| Box::pin(dns_delete(a, b, c)),
            },
            SubCommand {
                subcommand: "list",
                description: "List all DNS records.",
                nargs: 0,
                usage: "",
                op: |a, b, c| Box::pin(dns_list(a, b, c)),
            },
            SubCommand {
                subcommand: "create",
                description: "Create a DNS record with given domain and type.",
                nargs: 3,
                usage: "<domain> <type> <ip/host>",
                op: |a, b, c| Box::pin(dns_create(a, b, c)),
            },
        ],
    },
    Command {
        command: "light",
        short: "l",
        description: "Light operations.",
        subcommands: &[
            SubCommand {
                subcommand: "on",
                description: "Turn on the light.",
                nargs: 0,
                usage: "",
                op: |a, b, c| Box::pin(light_on(a, b, c)),
            },
            SubCommand {
                subcommand: "off",
                description: "Turn off the light.",
                nargs: 0,
                usage: "",
                op: |a, b, c| Box::pin(light_off(a, b, c)),
            },
            SubCommand {
                subcommand: "ww",
                description: "Set the light to warm white.",
                nargs: 0,
                usage: "",
                op: |a, b, c| Box::pin(light_ww(a, b, c)),
            },
            SubCommand {
                subcommand: "color",
                description: "Set the light to a color.",
                nargs: 3,
                usage: "<r> <g> <b>",
                op: |a, b, c| Box::pin(light_color(a, b, c)),
            },
            SubCommand {
                subcommand: "brightness",
                description: "Set the light brightness.",
                nargs: 1,
                usage: "<brightness>",
                op: |a, b, c| Box::pin(light_brightness(a, b, c)),
            },
        ],
    },
];

async fn create_connection() -> Result<(Channel, String), Box<dyn std::error::Error>> {
    let dirs = AppDirs::new(Some("moodyapi"), false).unwrap();
    let ini_path = dirs.config_dir.as_path().join("Client.ini");

    let conf = if let Ok(ini_file) = Ini::load_from_file(ini_path.as_path()) {
        ini_file
    } else if let Ok(ini_file) = Ini::load_from_file("/etc/moodyapi/Client.ini") {
        ini_file
    } else {
        println!("Failed to locate configurations.");
        exit(1);
    };

    let api_host = conf.general_section().get("Server").unwrap().to_string();
    let client_id = conf.general_section().get("ClientID").unwrap().to_string();

    let grpc_channel = Channel::from_shared(api_host.clone())?
        .connect()
        .await
        .expect("Can't create a channel");

    Ok((grpc_channel, client_id))
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    if env::args().len() < 2 || env::args().nth(1).unwrap() == "help" {
        println!("usage: {} <command> [args...]", env::args().nth(0).unwrap());
        println!("commands:");
        for cmd in COMMANDS {
            println!("  {} - {}", cmd.command, cmd.description);
        }

        exit(1);
    }

    let user_cmd = env::args().nth(1).unwrap().to_string();
    for cmd in COMMANDS {
        // if not match or not start with command, continue
        if cmd.command != user_cmd && cmd.short != user_cmd {
            continue;
        }

        if env::args().len() < 3 {
            println!(
                "usage: {} {} <subcommand> [args...]",
                env::args().nth(0).unwrap(),
                cmd.command
            );
            println!("subcommands:");
            for subcmd in cmd.subcommands {
                println!("  {} - {}", subcmd.subcommand, subcmd.description);
            }
            exit(1);
        }

        let user_subcmd = env::args().nth(2).unwrap().to_string();

        for subcmd in cmd.subcommands {
            if subcmd.subcommand != user_subcmd && !subcmd.subcommand.starts_with(&user_subcmd) {
                continue;
            }

            if env::args().len() < subcmd.nargs + 3 {
                println!(
                    "usage: {} {} {} {}",
                    env::args().nth(0).unwrap(),
                    cmd.command,
                    subcmd.subcommand,
                    subcmd.usage
                );
                exit(1);
            }

            let args = env::args().skip(3).collect::<Vec<String>>();
            let (channel, uuid) = create_connection().await?;
            match (subcmd.op)(channel, uuid, args).await {
                Ok(_) => exit(0),
                Err(e) => {
                    println!("Error: {}", e.message());
                    exit(1)
                }
            }
        }

        println!("Unknown subcommand: {}", user_subcmd);
        exit(1);
    }

    println!("Unknown command: {}", user_cmd);
    exit(1);
}

async fn notification_send(chan: Channel, uuid: String, args: Vec<String>) -> Result<(), Status> {
    println!("Sending notifications...");
    let channel_str = args[0].to_string();
    let n_channel: i64;

    if channel_str.ends_with("p") {
        n_channel = channel_str[..channel_str.len() - 1].parse().unwrap();
    } else {
        n_channel = channel_str.parse().unwrap();
    }

    let n = Notification {
        title: args[1].to_string(),
        content: args[2].to_string(),
        channel_id: n_channel,
        private: channel_str.ends_with("p"),
        ..Default::default()
    };

    MoodyApiServiceClient::new(chan)
        .send_notification(Request::new(SendRequest {
            auth: Some(Auth {
                client_uuid: uuid.to_owned(),
            }),
            notification: Some(n),
        }))
        .await
        .and_then(|_| Ok(()))
}

async fn notification_subscribe(
    chan: Channel,
    uuid: String,
    args: Vec<String>,
) -> Result<(), Status> {
    loop {
        let mut client = MoodyApiServiceClient::new(chan.clone());

        let request = Request::new(SubscribeRequest {
            auth: Some(Auth {
                client_uuid: uuid.clone().to_owned(),
            }),
            channel_id: args[0].parse().unwrap(),
            channels: vec![args[0].parse().unwrap()],
            urgency: Some(0), // todo urgency
        });

        match client.subscribe_notifications(request).await {
            Err(e) => println!("something went wrong: {}", e),
            Ok(stream) => {
                let mut resp_stream = stream.into_inner();
                loop {
                    match resp_stream.message().await {
                        Ok(None) => println!("expect a notification object"),
                        Ok(Some(n)) => {
                            println!("Received Notification: {:?}", n);
                            Notify::new()
                                .summary(&n.title)
                                .body(&n.content)
                                .icon(&n.icon)
                                .appname("Notify Client")
                                .hint(Hint::Resident(true))
                                .show()
                                .unwrap();
                        }
                        Err(e) => {
                            println!("something went wrong: {}", &e);
                            break;
                        }
                    }
                    sleep(Duration::from_secs(2)).await;
                }
            }
        }

        sleep(Duration::from_secs(10)).await;
    }
}

async fn channel_create(chan: Channel, uuid: String, args: Vec<String>) -> Result<(), Status> {
    let mut client = MoodyApiServiceClient::new(chan);

    let request = Request::new(CreateChannelRequest {
        auth: Some(Auth {
            client_uuid: uuid.to_owned(),
        }),
        channel: Some(NotificationChannel {
            name: args[0].to_owned(),
            ..Default::default()
        }),
    });

    match client.create_notification_channel(request).await {
        Err(e) => println!("something went wrong: {}", e),
        Ok(resp) => {
            let resp = resp.into_inner();
            println!("Channel created: {:?}", resp.channel);
        }
    }

    Ok(())
}

async fn channel_delete(chan: Channel, uuid: String, args: Vec<String>) -> Result<(), Status> {
    let mut client = MoodyApiServiceClient::new(chan);

    let request = Request::new(DeleteChannelRequest {
        auth: Some(Auth {
            client_uuid: uuid.to_owned(),
        }),
        channel_id: args[0].parse().unwrap(),
    });

    client
        .delete_notification_channel(request)
        .await
        .and_then(|_| {
            println!("Channel deleted");
            Ok(())
        })
}

async fn channel_list(chan: Channel, uuid: String, _args: Vec<String>) -> Result<(), Status> {
    let mut client = MoodyApiServiceClient::new(chan);

    let request = Request::new(ListChannelRequest {
        auth: Some(Auth {
            client_uuid: uuid.to_owned(),
        }),
    });

    client
        .list_notification_channel(request)
        .await
        .and_then(|resp| {
            let resp = resp.into_inner();
            for channel in resp.channels {
                println!("Channel {}: {}", channel.id, channel.name);
            }
            Ok(())
        })
}

async fn _do_camera_state(chan: Channel, uuid: String, state: bool) -> Result<(), Status> {
    let mut client = MoodyApiServiceClient::new(chan);

    let request = Request::new(UpdateCameraStateRequest {
        auth: Some(Auth {
            client_uuid: uuid.to_owned(),
        }),
        state: Some(CameraState { state }),
    });

    client.set_camera_state(request).await.and_then(|_| {
        println!("Camera is now {}", if state { "on" } else { "off" });
        Ok(())
    })
}

async fn camera_start(chan: Channel, uuid: String, _args: Vec<String>) -> Result<(), Status> {
    _do_camera_state(chan, uuid, true).await
}

async fn camera_stop(chan: Channel, uuid: String, _args: Vec<String>) -> Result<(), Status> {
    _do_camera_state(chan, uuid, false).await
}

async fn dns_update(chan: Channel, uuid: String, args: Vec<String>) -> Result<(), Status> {
    let mut client = MoodyApiServiceClient::new(chan);

    let request = Request::new(UpdateDnsRecordRequest {
        auth: Some(Auth {
            client_uuid: uuid.to_owned(),
        }),
        record: Some(DnsRecord {
            name: args[0].to_owned(),
            r#type: args[1].to_owned(),
            ip: args[2].to_owned(),
        }),
    });

    client.update_dns_record(request).await.and_then(|_| {
        println!("DNS updated");
        Ok(())
    })
}

async fn dns_delete(chan: Channel, uuid: String, args: Vec<String>) -> Result<(), Status> {
    let mut client = MoodyApiServiceClient::new(chan);

    let request = Request::new(DeleteDnsRecordRequest {
        auth: Some(Auth {
            client_uuid: uuid.to_owned(),
        }),
        name: args[0].to_owned(),
        r#type: args[1].to_owned(),
    });

    client.delete_dns_record(request).await.and_then(|_| {
        println!("DNS deleted");
        Ok(())
    })
}

async fn dns_list(chan: Channel, uuid: String, _args: Vec<String>) -> Result<(), Status> {
    let mut client = MoodyApiServiceClient::new(chan);

    let request = Request::new(ListDnsRecordsRequest {
        auth: Some(Auth {
            client_uuid: uuid.to_owned(),
        }),
        ..Default::default()
    });

    client.list_dns_records(request).await.and_then(|resp| {
        let resp = resp.into_inner();
        for record in resp.entries {
            println!("{} {} {}", record.name, record.r#type, record.ip);
        }
        Ok(())
    })
}

async fn dns_create(chan: Channel, uuid: String, args: Vec<String>) -> Result<(), Status> {
    let mut client = MoodyApiServiceClient::new(chan);

    let request = Request::new(CreateDnsRecordRequest {
        auth: Some(Auth {
            client_uuid: uuid.to_owned(),
        }),
        record: Some(DnsRecord {
            name: args[0].to_owned(),
            r#type: args[1].to_owned(),
            ip: args[2].to_owned(),
        }),
    });

    client.create_dns_record(request).await.and_then(|_| {
        println!("DNS created");
        Ok(())
    })
}

async fn light_on(chan: Channel, uuid: String, _args: Vec<String>) -> Result<(), Status> {
    let mut client = MoodyApiServiceClient::new(chan);

    let request = Request::new(SetLightRequest {
        auth: Some(Auth {
            client_uuid: uuid.to_owned(),
        }),
        state: Some(LightState {
            on: true,
            ..Default::default()
        }),
    });

    client.set_light_state(request).await.and_then(|_| {
        println!("Light is now on");
        Ok(())
    })
}

async fn light_off(chan: Channel, uuid: String, _args: Vec<String>) -> Result<(), Status> {
    let mut client = MoodyApiServiceClient::new(chan);

    let request = Request::new(SetLightRequest {
        auth: Some(Auth {
            client_uuid: uuid.to_owned(),
        }),
        state: Some(LightState {
            on: false,
            ..Default::default()
        }),
    });

    client.set_light_state(request).await.and_then(|_| {
        println!("Light is now off");
        Ok(())
    })
}

async fn light_color(chan: Channel, uuid: String, args: Vec<String>) -> Result<(), Status> {
    let mut client = MoodyApiServiceClient::new(chan);

    let request = Request::new(SetLightRequest {
        auth: Some(Auth {
            client_uuid: uuid.to_owned(),
        }),
        state: Some(LightState {
            on: true,
            brightness: 255,
            mode: Some(Mode::Colored(LightColor {
                red: args[0].parse().unwrap(),
                green: args[1].parse().unwrap(),
                blue: args[2].parse().unwrap(),
            })),
        }),
    });

    client.set_light_state(request).await.and_then(|_| {
        println!("Light color set");
        Ok(())
    })
}

async fn light_ww(chan: Channel, uuid: String, _args: Vec<String>) -> Result<(), Status> {
    let mut client = MoodyApiServiceClient::new(chan);

    let request = Request::new(SetLightRequest {
        auth: Some(Auth {
            client_uuid: uuid.to_owned(),
        }),
        state: Some(LightState {
            on: true,
            brightness: 255,
            mode: Some(Mode::Warmwhite(true)),
        }),
    });

    client.set_light_state(request).await.and_then(|_| {
        println!("Light color set");
        Ok(())
    })
}

async fn light_brightness(chan: Channel, uuid: String, args: Vec<String>) -> Result<(), Status> {
    let mut client = MoodyApiServiceClient::new(chan);

    let brightness: u32 = if args[0] == "max" {
        127
    } else if args[0] == "min" {
        5
    } else {
        args[0].parse::<u32>().unwrap().min(127).max(5)
    };

    let request = Request::new(SetLightRequest {
        auth: Some(Auth {
            client_uuid: uuid.to_owned(),
        }),
        state: Some(LightState {
            on: true,
            brightness,
            ..Default::default()
        }),
    });

    client.set_light_state(request).await.and_then(|_| {
        println!("Light brightness set");
        Ok(())
    })
}
