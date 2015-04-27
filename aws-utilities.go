package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "aws-utilities"
	app.Version = Version
	app.Usage = ""
	app.Author = "simeji"
	app.Email = "simeji.net@gmail.com"
	app.Commands = Commands
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "profile, p",
			Usage: "[*required] aws profile",
		},
	}
	app.Before = func(c *cli.Context) (err error) {
		err = nil
		if c.GlobalString("profile") == "" {
			err = fmt.Errorf("'--profile' is required")
			fmt.Println(err)
			os.Exit(1)
		}
		return
	}
	app.Run(os.Args)
}
