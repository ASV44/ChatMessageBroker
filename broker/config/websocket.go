package config

import (
	"github.com/mitchellh/mapstructure"
	"time"
)

// WebsocketConnectionSettings contains websocket related config parameters
type WebsocketConnectionSettings struct {
	ReadBufferSize  int   `mapstructure:"read_buffer_size"`
	WriteBufferSize int   `mapstructure:"write_buffer_size"`
	MaxMessageSize  int64 `mapstructure:"max_message_size"`
	WriteWait       time.Duration
	PongWait        time.Duration
	PingPeriod      time.Duration
}

// NewWebsocketConnectionSettings decodes websocket config value from map and returns new instance of WebsocketConnectionSettings
func NewWebsocketConnectionSettings(configManager Manager) (WebsocketConnectionSettings, error) {
	var websocketSettings WebsocketConnectionSettings
	err := mapstructure.Decode(configManager.websocketMapConfig(), &websocketSettings)
	if err != nil {
		return websocketSettings, err
	}

	websocketSettings.WriteWait = configManager.websocketWriteWait()
	websocketSettings.PongWait = configManager.websocketPongWait()
	websocketSettings.PingPeriod = configManager.websocketPingPeriod()

	return websocketSettings, nil
}
