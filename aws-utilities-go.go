package main

import (
  "os"
  "github.com/codegangsta/cli"
)

func main() {
  app := cli.NewApp()
  app.Name = "aws-utilities"
  app.Version = Version
  app.Usage = ""
  app.Author = "simeji"
  app.Email = "simeji.net@gmail.com"
  app.Commands = Commands

  app.Run(os.Args)
}
