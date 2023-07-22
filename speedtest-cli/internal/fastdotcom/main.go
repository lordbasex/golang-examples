package fastdotcom

import (
	"context"
	"log"

	"go.jonnrb.io/speedtest/fastdotcom"
)

func Main(args []string) {
	err := flagSet.Parse(args[1:])
	if err != nil {
		panic(err)
	}

	var client fastdotcom.Client

	ctx, cancel := context.WithTimeout(context.Background(), *cfgTime)
	defer cancel()

	m, err := fastdotcom.GetManifest(ctx, *urlCount)
	if err != nil {
		log.Fatalf("Error loading fast.com configuration: %v", err)
	}

	download(m, &client)
	upload(m, &client)
}
