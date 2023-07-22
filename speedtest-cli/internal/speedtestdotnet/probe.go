package speedtestdotnet

import (
	"context"
	"fmt"
	"log"

	"go.jonnrb.io/speedtest/oututil"
	"go.jonnrb.io/speedtest/speedtestdotnet"
	"go.jonnrb.io/speedtest/units"
	"golang.org/x/sync/errgroup"
)

func download(client *speedtestdotnet.Client, server speedtestdotnet.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), *dlTime)
	defer cancel()

	stream, finalize := proberPrinter(func(s units.BytesPerSecond) string {
		return formatSpeed("Download speed", s)
	})
	speed, err := server.ProbeDownloadSpeed(ctx, client, stream)
	if err != nil {
		log.Fatalf("Error probing download speed: %v", err)
		return
	}
	finalize(speed)
}

func upload(client *speedtestdotnet.Client, server speedtestdotnet.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), *ulTime)
	defer cancel()

	stream, finalize := proberPrinter(func(s units.BytesPerSecond) string {
		return formatSpeed("Upload speed", s)
	})
	speed, err := server.ProbeUploadSpeed(ctx, client, stream)
	if err != nil {
		log.Fatalf("Error probing upload speed: %v", err)
	}
	finalize(speed)
}

func proberPrinter(format func(units.BytesPerSecond) string) (
	stream chan units.BytesPerSecond,
	finalize func(units.BytesPerSecond),
) {
	p := oututil.StartPrinting()
	p.Println(format(units.BytesPerSecond(0)))

	stream = make(chan units.BytesPerSecond)
	var g errgroup.Group
	g.Go(func() error {
		for speed := range stream {
			p.Println(format(speed))
		}
		return nil
	})

	finalize = func(s units.BytesPerSecond) {
		g.Wait()
		p.Finalize(format(s))
	}
	return
}

func formatSpeed(prefix string, s units.BytesPerSecond) string {
	var i interface{}
	// Default return speed is in bytes.
	if *fmtBytes {
		i = s
	} else {
		i = s.BitsPerSecond()
	}
	return fmt.Sprintf("%s: %v", prefix, i)
}
