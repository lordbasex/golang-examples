package fastdotcom

import (
	"flag"
	"time"
)

var (
	flagSet = flag.NewFlagSet("fastdotcom", flag.ExitOnError)

	fmtBytes = flagSet.Bool("bytes", false, "Display speeds in SI bytes (default is bits)")
	urlCount = flagSet.Int("urls", 5, "Number of URLs to use to probe")
	cfgTime  = flagSet.Duration("time.config", 1*time.Second, "Timeout for getting initial configuration")
	pngTime  = flagSet.Duration("time.latency", 1*time.Second, "Timeout for latency detection phase")
	dlTime   = flagSet.Duration("time.download", 10*time.Second, "Maximum time to spend in download probe phase")
	ulTime   = flagSet.Duration("time.upload", 10*time.Second, "Maximum time to spend in upload probe phase")
)
