package cmd

import (
	"bufio"
	"fmt"
	"github.com/ksp237/peerbeam/conn"
	"github.com/ksp237/peerbeam/receiver"
	"github.com/ksp237/peerbeam/sender"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func App() error {
	var app = &cli.App{
		Name:                 "peerbeam",
		Usage:                "A simple and secure p2p file transfer tool powered by WebRTC",
		EnableBashCompletion: true,
		//Action: func(c *cli.Context) error {
		//	return nil
		//},
		Commands: []*cli.Command{
			{
				Name:      "send",
				Usage:     "Send files",
				UsageText: "peerbeam send [files...]",
				Action: func(c *cli.Context) error {
					files := c.Args().Slice()
					fmt.Println("files:", files)
					if len(files) == 0 {
						fmt.Print("Enter the file(s) to send (space seperated): ")
						reader := bufio.NewReader(os.Stdin)
						input, err := reader.ReadString('\n')
						if err != nil {
							return err
						}
						files = strings.Split(strings.TrimSpace(input), ",")
					}

					err := sender.StartSender(files)
					if err != nil {
						return err
					}

					return nil
				},
			},
			{
				Name:      "receive",
				Usage:     "Receive files",
				UsageText: "peerbeam receive [destination]",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:      "dest",
						Aliases:   []string{"d"},
						Usage:     "Path to destination directory",
						Required:  false,
						TakesFile: true,
					},
				},
				Action: func(c *cli.Context) error {
					dest := c.String("dest")
					dest = strings.TrimSpace(dest)
					if dest == "" {
						fmt.Print("Enter the destination path: ")
						reader := bufio.NewReader(os.Stdin)
						input, err := reader.ReadString('\n')
						if err != nil {
							return err
						}
						dest = strings.TrimSpace(input)
						if dest == "" {
							dest = "."
						}
					}
					destPath, err := filepath.Abs(dest)
					if err != nil {
						return fmt.Errorf("error with destination path '%s': %v", dest, err)
					}

					destInfo, err := os.Stat(destPath)
					if err != nil {
						return fmt.Errorf(err.Error())
					}
					if !destInfo.IsDir() {
						return fmt.Errorf("destination path must be a directory")
					}

					fmt.Printf("Receiving file and saving to %s\n", destPath)

					err = receiver.StartReceiver(destPath)
					if err != nil {
						return err
					}

					return nil
				},
			},
			{
				Name:      "stun",
				Usage:     "Query STUN server for connection info (server-reflexive ICE candidates)",
				UsageText: "peerbeam stun [flags]",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:     "ip-only",
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
