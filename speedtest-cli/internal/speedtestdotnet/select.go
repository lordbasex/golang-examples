package speedtestdotnet

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.jonnrb.io/speedtest/geo"
	"go.jonnrb.io/speedtest/speedtestdotnet"
)

// Selects a server to use, either selected by the user or by a low latency
// selection algorithm.
//
func selectServer(client *speedtestdotnet.Client, cfg speedtestdotnet.Config, servers []speedtestdotnet.Server) speedtestdotnet.Server {
	var (
		distance geo.Kilometers
		latency  time.Duration
		server   speedtestdotnet.Server
	)

	ctx, cancel := context.WithTimeout(context.Background(), *pngTime)
	defer cancel()

	if *srvID != 0 {
		id := speedtestdotnet.ServerID(*srvID)

		// Meh, linear search.
		i := -1
		for j, s := range servers {
			if s.ID == id {
				i = j
				break
			}
		}
		if i == -1 {
			log.Fatalf("Server not found: %d\n", id)
		}

		server = servers[i]
		l, err := server.AverageLatency(ctx, client, speedtestdotnet.DefaultLatencySamples)
		if err != nil {
			log.Fatalf("Error getting latency for (%v): %v", server, err)
		}

		latency = l
		distance = cfg.Coordinates.DistanceTo(server.Coordinates)
	} else {
		distanceMap := speedtestdotnet.SortServersByDistance(servers, cfg.Coordinates)

		// Truncate to just a few of the closest servers for the latency test.
		const maxCloseServers = 5
		closestServers := func() []speedtestdotnet.Server {
			if len(servers) > maxCloseServers {
				return servers[:maxCloseServers]
			} else {
				return servers
			}
		}()

		latencyMap, err := speedtestdotnet.StableSortServersByAverageLatency(
			closestServers, ctx, client, speedtestdotnet.DefaultLatencySamples)
		if err != nil {
			log.Fatalf("Error getting server latencies: %v", err)
		}

		server = closestServers[0]
		latency = latencyMap[server.ID]
		distance = distanceMap[server.ID]
	}

	fmt.Printf("Using server %d hosted by %s (%s) [%v]: %.1f ms\n",
		server.ID, server.Sponsor, server.Name, distance, float64(latency)/float64(time.Millisecond))

	return server
}
