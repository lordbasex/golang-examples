package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"speedtest-cli/internal/fastdotcom"
	"speedtest-cli/internal/speedtestdotnet"
)

type subcmd struct {
	mainFunc func(args []string)
	aliases  []string
}

var subcmds = []subcmd{
	subcmd{
		mainFunc: speedtestdotnet.Main,
		aliases:  []string{"st", "speedtest.net"},
	},
	subcmd{
		mainFunc: fastdotcom.Main,
		aliases:  []string{"f", "fast.com"},
	},
}

func main() {
	flag.Usage = usage
	flag.Parse()

	s := getSubcmd()
	if s == nil {
		flag.Usage()
		os.Exit(2)
	}
	s.mainFunc(flag.Args())
}

func getSubcmd() *subcmd {
	args := flag.Args()
	if len(args) < 1 {
		return nil
	}
	for _, s := range subcmds {
		for _, a := range s.aliases {
			if a == args[0] {
				return &s
			}
		}
	}
	return nil
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "USAGE\n")
	for _, s := range subcmds {
		fmt.Fprintf(
			flag.CommandLine.Output(),
			"  %s %s [OPTIONS]\n",
			os.Args[0], strings.Join(s.aliases, "|"))
	}
	flag.PrintDefaults()
}
