package cmd

import (
	"fmt"
	"github.com/6b70/peerbeam/conn"
	"github.com/urfave/cli/v2"
	"os"
	"strconv"
	"strings"
)

func App() error {
	var app = &cli.App{
		Name:                 "peerbeam",
		Usage:                "A simple and secure p2p file transfer tool powered by WebRTC",
		EnableBashCompletion: true,
		UsageText:            "peerbeam <command> [options] [file(s)/destination]",
		Commands: []*cli.Command{
			{
				Name:      "send",
				Usage:     "Send files",
				UsageText: "peerbeam send <filename>...",
				Action: func(c *cli.Context) error {
					files := c.Args().Slice()
					if len(files) == 0 {
						return fmt.Errorf("no files specified: use 'peerbeam send <filename>...'")
					}
					return startSender(files)
				},
			},
			{
				Name:      "receive",
				Aliases:   []string{"recv"},
				Usage:     "Receive files",
				UsageText: "peerbeam receive [destination]",
				Action: func(c *cli.Context) error {
					recvArgs := c.Args().Slice()
					if len(recvArgs) > 1 {
						return fmt.Errorf("too many arguments for receive: use 'peerbeam receive [destination]'")
					}
					var dest string
					if len(recvArgs) == 1 {
						dest = recvArgs[0]
					} else {
						dest = "."
					}
					dest = strings.TrimSpace(dest)
					return startReceiver(dest)
				},
			},
			{
				Name:      "stun",
				Usage:     "Query STUN server for connection info (server-reflexive ICE candidates)",
				UsageText: "peerbeam stun [flags]",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:     "ip",
						Aliases:  []string{"i"},
						Usage:    "Only display ip:port pairs",
						Required: false,
					},
				},
				Action: func(c *cli.Context) error {
					res, err := conn.FetchSRFLX()
					if err != nil {
						return err
					}
					if len(res) == 0 {
						return fmt.Errorf("no srflx candidates found")
					}
					ipOnly := c.Bool("ip-only")
					for _, ic := range res {
						if ipOnly {
							fmt.Println(ic.Address + ":" + strconv.Itoa(int(ic.Port)))
						} else {
							fmt.Println(ic.String())
						}
					}
					return nil
				},
			},
		},
	}
	app.Setup()
	if err := app.Run(os.Args); err != nil {
		return err
	}
	return nil
}
