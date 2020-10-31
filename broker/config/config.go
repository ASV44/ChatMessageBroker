package config

import (
	"github.com/spf13/viper"
	"time"
)

// config yaml paths
const (
	workspace = "workspace"

	tcpAddress        = "tcp_server.address"
	tcpConnectionType = "tcp_server.connection_type"

	httpAddress = "http_server.address"

	shutdownTimeout = "timeout.shutdown"
	readTimeout     = "timeout.read"
	writeTimeout    = "timeout.write"

	websocketReadBufferSize  = "websocket.read_buffer_size"
	websocketWriteBufferSize = "websocket.write_buffer_size"
	websocketMaxMessageSize  = "websocket.max_message_size"
	websocketWriteWait       = "websocket.write_wait"
	websocketPongWait        = "websocket.pong_wait"
	websocketPingPeriod      = "websocket.ping_period"

	logLevel = "logging.level"

	httpTimeout = "http.conn_timeout"
)

// Manager is holder of configuration related data
type Manager struct {
	Path        string
	viperConfig *viper.Viper
}

// NewManager initializes configs
func NewManager(configPath string) (Manager, error) {
	configManager := &Manager{Path: configPath, viperConfig: viper.New()}
	configManager.SetDefaults()

	err := configManager.readFromFile(configPath)
	if err != nil {
		return *configManager, err
	}

	return *configManager, nil
}

// SetDefaults sets default values for all configuration parameters
func (manager *Manager) SetDefaults() {
	manager.viperConfig.SetDefault(workspace, "No name")

	manager.viperConfig.SetDefault(tcpAddress, "localhost:8888")
	manager.viperConfig.SetDefault(tcpConnectionType, "tcp")

	manager.viperConfig.SetDefault(httpAddress, ":8080")

	manager.viperConfig.SetDefault(shutdownTimeout, 15*time.Second)
	manager.viperConfig.SetDefault(readTimeout, 10*time.Second)
	manager.viperConfig.SetDefault(writeTimeout, 10*time.Second)

	manager.viperConfig.SetDefault(websocketReadBufferSize, 1024)
	manager.viperConfig.SetDefault(websocketWriteBufferSize, 1024)
	manager.viperConfig.SetDefault(websocketMaxMessageSize, 512)
	manager.viperConfig.SetDefault(websocketWriteWait, 10*time.Second)
	manager.viperConfig.SetDefault(websocketPongWait, 60*time.Second)
	manager.viperConfig.SetDefault(websocketPingPeriod, 60*0.9*time.Second)

	manager.viperConfig.SetDefault(httpTimeout, 1*time.Second)

	manager.viperConfig.SetDefault(logLevel, "debug")
}

func (manager *Manager) readFromFile(filename string) error {
	manager.viperConfig.SetConfigFile(filename)
	return manager.viperConfig.ReadInConfig()
}

// Workspace returns workspace name
func (manager Manager) Workspace() string {
	return manager.viperConfig.GetString(workspace)
}

// GetLoggingLevel returns config settings for logger
func (manager Manager) GetLoggingLevel() string {
	return manager.viperConfig.GetString(logLevel)
}

// TCPAddress returns tcp services address
func (manager Manager) TCPAddress() string {
	return manager.viperConfig.GetString(tcpAddress)
}

// TCPServerConnectionType returns connection type used for net.Listen
func (manager Manager) TCPServerConnectionType() string {
	return manager.viperConfig.GetString(tcpConnectionType)
}

// HTTPAddress returns port number for HTTP server
func (manager Manager) HTTPAddress() string {
	return manager.viperConfig.GetString(httpAddress)
}

// ReadTimeout returns read timeout for http server
func (manager Manager) ReadTimeout() time.Duration {
	return manager.viperConfig.GetDuration(readTimeout)
}

// WriteTimeout returns write timeout for http server
func (manager Manager) WriteTimeout() time.Duration {
	return manager.viperConfig.GetDuration(writeTimeout)
}

// websocketMapConfig returns websocket map configuration
func (manager Manager) websocketMapConfig() map[string]interface{} {
	return manager.viperConfig.GetStringMap("websocket")
}

// websocketWriteWait returns write wait for websocket connection
func (manager Manager) websocketWriteWait() time.Duration {
	return manager.viperConfig.GetDuration(websocketWriteWait)
}

// websocketPongWait returns pong wait duration for websocket connection
func (manager Manager) websocketPongWait() time.Duration {
	return manager.viperConfig.GetDuration(websocketPongWait)
}

// websocketPingPeriod returns ping period duration for websocket connection
func (manager Manager) websocketPingPeriod() time.Duration {
	pingPeriod := manager.viperConfig.GetDuration(websocketPingPeriod)
	pongWait := manager.websocketPongWait()
	if pingPeriod > pongWait {
		pingPeriod = time.Duration(pongWait.Seconds() * 0.9)
	}

	return pingPeriod
}

// HTTPTimeout returns timeout to be used by default in external HTTP requests
func (manager Manager) HTTPTimeout() time.Duration {
	return manager.viperConfig.GetDuration(httpTimeout)
}
