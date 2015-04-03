package main

import (
  "log"
  "os"
  "github.com/codegangsta/cli"
  //"github.com/awslabs/aws-sdk-go/aws"
  //"github.com/awslabs/aws-sdk-go/service/ec2"
)

var Commands = []cli.Command{
  commandList_ipaddress,
  commandEnter_instance,
}


var commandList_ipaddress = cli.Command{
  Name:  "list-ipaddress",
  Aliases:  []string{"li"},
  Usage: "get ipaddress list by NameTag",
  Description: `
  hogehoge mogemoge
`,
  Action: doList_ipaddress,
}

var commandEnter_instance = cli.Command{
  Name:  "enter-instance",
  Aliases:  []string{"ei"},
  Usage: "",
  Description: `
  aaaa hogehoge mogemoge
`,
  Action: doEnter_instance,
}


func debug(v ...interface{}) {
  if os.Getenv("DEBUG") != "" {
    log.Println(v...)
  }
}

func assert(err error) {
  if err != nil {
    log.Fatal(err)
  }
}


func doList_ipaddress(c *cli.Context) {
  name := "default"
  if len(c.Args()) > 0 {
    name = c.Args()[0]
  } else {
    name = c.String("profile")
  }
  log.Println(os.Getenv("USER") + "||" + name)
}

func doEnter_instance(c *cli.Context) {
}


