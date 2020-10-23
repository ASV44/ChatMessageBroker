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

	shutdownTimeout = "timeout.shutdown"
	readTimeout     = "timeout.read"
	writeTimeout    = "timeout.write"

	logLevel = "logging.level"

	httpTimeout = "http.conn_timeout"
)

// Manager is holder of configuration related data
type Manager struct {
	Path        string
	viperConfig *viper.Viper
}

// NewManager initializes configs
func NewManager(configPath string) (*Manager, error) {
	configManager := &Manager{Path: configPath, viperConfig: viper.New()}
	configManager.SetDefaults()

	err := configManager.readFromFile(configPath)
	if err != nil {
		return nil, err
	}

	return configManager, nil
}

// SetDefaults sets default values for all configuration parameters
func (manager *Manager) SetDefaults() {
	manager.viperConfig.SetDefault(workspace, "No name")

	manager.viperConfig.SetDefault(tcpAddress, "localhost:8888")
	manager.viperConfig.SetDefault(tcpConnectionType, "tcp")

	manager.viperConfig.SetDefault(shutdownTimeout, 15*time.Second)
	manager.viperConfig.SetDefault(readTimeout, 10*time.Second)
	manager.viperConfig.SetDefault(writeTimeout, 10*time.Second)

	manager.viperConfig.SetDefault(httpTimeout, 1*time.Second)

	manager.viperConfig.SetDefault(logLevel, "debug")
}

func (manager *Manager) readFromFile(filename string) error {
	manager.viperConfig.SetConfigFile(filename)
	return manager.viperConfig.ReadInConfig()
}

// Workspace returns workspace name
func (manager *Manager) Workspace() string {
	return manager.viperConfig.GetString(workspace)
}

// GetLoggerSettings returns config settings for logger
func (manager *Manager) GetLoggerSettings() map[string]interface{} {
	return manager.viperConfig.GetStringMap("logging")
}

// TCPAddress returns tcp services address
func (manager *Manager) TCPAddress() string {
	return manager.viperConfig.GetString(tcpAddress)
}

// TCPServerConnectionType returns connection type used for net.Listen
func (manager *Manager) TCPServerConnectionType() string {
	return manager.viperConfig.GetString(tcpConnectionType)
}

// ReadTimeout returns read timeout for http server
func (manager *Manager) ReadTimeout() time.Duration {
	return manager.viperConfig.GetDuration(readTimeout)
}

// WriteTimeout returns write timeout for http server
func (manager *Manager) WriteTimeout() time.Duration {
	return manager.viperConfig.GetDuration(writeTimeout)
}

// HTTPTimeout returns timeout to be used by default in external HTTP requests
func (manager *Manager) HTTPTimeout() time.Duration {
	return manager.viperConfig.GetDuration(httpTimeout)
}
