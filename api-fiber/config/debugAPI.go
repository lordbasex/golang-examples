package config

import (
	"log"
	"strconv"
)

// DebugConfiguration represents the debug mode status.
var DebugConfiguration bool = false

// DebugAPI sets up the debug mode for the API.
func DebugAPI() {
	// Set log flags to include timestamps and short file names
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Retrieve the DEBUG_CONFIGURATION value from the environment
	s := Config("DEBUG_CONFIGURATION")

	// Parse the value as a boolean and update DebugConfiguration
	DebugConfiguration, _ = strconv.ParseBool(s)
}
