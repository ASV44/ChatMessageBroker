package models

import "fmt"

// PluginError represents abstraction for server error
type PluginError struct {
	Message  string
	Err      error
}

func (e PluginError) Error() string {
	return fmt.Sprintf("Plugin Error : %s \nCaused by : %s", e.Message, e.Err)
}
