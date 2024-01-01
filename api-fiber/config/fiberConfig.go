package config

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// FiberListen represents the server URL configured from environment variables.
var FiberListen string = Config("SERVER_URL")

// FiberConfig returns the configuration for the Fiber web framework.
// See: https://docs.gofiber.io/api/fiber#config
func FiberConfig() fiber.Config {
	// Define server settings.
	var readTimeoutSecondsCount int
	var err error

	// Parse the read timeout from the environment variable.
	readTimeoutSecondsCount, err = strconv.Atoi(Config("SERVER_READ_TIMEOUT"))
	if err != nil {
		// Set a default value if parsing fails.
		readTimeoutSecondsCount = 60
	}

	// Configure and return the Fiber configuration.
	return fiber.Config{
		AppName:               "API Fiber v1.0.0",
		ServerHeader:          "MyApiFiberDemo",
		DisableStartupMessage: false,                                                // Disable the startup banner displayed when the application starts. Default is false.
		Prefork:               true,                                                 // Enable or disable the prefork option. If enabled, Fiber will use prefork functionality if possible. Default is true.
		ReadTimeout:           time.Second * time.Duration(readTimeoutSecondsCount), // The amount of time allowed to read the full request, including the body. The default timeout is unlimited.
	}
}
