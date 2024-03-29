# Chat Message Broker
Golang Message Broker of Chat application, part of distributed system

Project includes both server (broker) and client apps.

## Install

Clone repo in to your machine and run `go mod download`

### Broker
Change directory to broker folder which is the main directory for broker app.

Check config yaml file for adopting all available broker configuration.

### Build
You should have working setup of GO  environment.

Run `make` in the project root. Start `broker_app` and `client_app`.

### Run
#### Broker
Broker app accepts: 
- `--config` flag which represents the path to config file. <br> 
Default value for `--config` flag is `./broker/config.yaml`. Modify path to config file only in case you use
or change config path location.
- `--connection` flag which represents socket connection type.<br> Broker supports 
`tcp-socket` and `websocket` connections. Default flag value is `all` which enables
all connection types. 

#### Client
Client app accepts:
- `--host` flag which represents Broker host address. Default value is `0.0.0.0`
- `--port` flag which represents Broker host address port number. Default value is `8888`
- `--connection-type` flag which represents Broker connection. Default value is `tcp`

Pass new values only in case you want to modify default flag values.

## How to use
Broker accepts 3 types of messages
- Direct message
- Channel message
- Command message

Client app defines for each type of message special symbol
- Direct message -> `@`. Message example: `@{$USER_NAME} {$MESSAGE_TEXT}`
- Channel message -> `#`. Message example: `#{$CHANNEL_NAME} {$MESSAGE_TEXT}`
- Command message -> `/`. Message example: `/{$SUPPORTED_COMMAND}`

At broker app start is created one default channel `random`. 

At client connection, broker app registers new user. All users are by default subscribed
to `random` channel.

Broker supports various commands which allow users to interact with broker workspace.

Broker supported commands:
- `create`- creates a new channel. User who created channel is subscribed to it by default.
 Command example: `/create {$CHANNEL_NAME}`
- `join` - joins a specific channel. Command example: `/join {$CHANNEL_NAME}`
- `leave` - joins a specific channel. Command example: `/leave {$CHANNEL_NAME}`
- `show` - displays specific info about workspace or channel. Command example: `/show {$OPTION}`
    Show command supports 4 options: 
    - `users` - displays list of all workspace users.
    - `channels` - display list of all workspace channels
    - `all` - display list of all workspace users and channels
    - `$CHANNEL_NAME` - display list of all users subscribed to a specific channel. Ex. `/show random` 
