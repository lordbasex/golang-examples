package speedtestdotnet

import (
	"context"
	"fmt"
	"log"

	"go.jonnrb.io/speedtest/speedtestdotnet"
)

func Main(args []string) {
	err := flagSet.Parse(args[1:])
	if err != nil {
		panic(err)
	}

	var client speedtestdotnet.Client

	if *list {
		printServers(&client)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), *cfgTime)
	defer cancel()

	cfg, err := client.Config(ctx)
	if err != nil {
		log.Fatalf("Error loading speedtest.net configuration: %v", err)
	}
	fmt.Printf("Testing from %s (%s)...\n", cfg.ISP, cfg.IP)
	servers := listServers(ctx, &client)

	server := selectServer(&client, cfg, servers)

	download(&client, server)
	upload(&client, server)
}
