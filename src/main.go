package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

var (
	PORT       = 53535
	CLIENT_MSG = []byte("pfg_ip_broadcast_cl")
	SERVER_MSG = []byte("pfg_ip_response_serv")
)

func main() {
	app := &cli.App{
		Name:                 "whichip",
		Usage:                "discover device's IP(s) in your local network",
		EnableBashCompletion: true,
		Commands: []cli.Command{
			{
				Name:  "version",
				Usage: "get version",
				Action: func(c *cli.Context) {
					fmt.Printf("%s\n", VERSION)
				},
			},
			{
				Name:  "listen",
				Usage: "listen to udp broadcast and respond accordingly",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "debug", Usage: "print debug message"},
					&cli.IntFlag{Name: "port", Usage: "the port to bind", Value: PORT},
				},
				Action: Listen,
			},
			{
				Name:  "discover",
				Usage: "make udp broadcast to discover device's IP(s)",
				Flags: []cli.Flag{
					&cli.Float64Flag{Name: "timeout", Usage: "discover timeout", Value: 1.0},
					&cli.BoolFlag{Name: "all", Usage: "print all IPs"},
					&cli.BoolFlag{Name: "debug", Usage: "print debug message"},
					&cli.IntFlag{Name: "port", Usage: "the port to bind", Value: PORT},
				},
				Action: Discover,
			},
		},
		Flags: []cli.Flag{
			&cli.Float64Flag{Name: "timeout", Usage: "discover timeout", Value: 1.0},
			&cli.BoolFlag{Name: "all", Usage: "print all IPs"},
			&cli.BoolFlag{Name: "debug", Usage: "print debug message"},
			&cli.IntFlag{Name: "port", Usage: "the port to bind", Value: PORT},
		},
		Action: Discover,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
