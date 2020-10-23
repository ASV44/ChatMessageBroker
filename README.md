# Chat Message Broker
Golang Message Broker of Chat application, part of distributed system

Project includes both server (broker) and client apps.

## Install

Clone repo in `$GOPATH`/src/github.com/ASV44/ folder.

### Broker
Change directory to broker folder which is the main directory for broker app.

Run `dep ensure`

Check config yaml file in for adopting all available broker configuration.

Run `make` in the project root. Start `broker_app` and `client_app`.
