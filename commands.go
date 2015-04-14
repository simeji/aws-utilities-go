package main

import (
  "log"
  "fmt"
  "os"
  "time"
  "github.com/codegangsta/cli"
  "github.com/awslabs/aws-sdk-go/aws"
  "github.com/awslabs/aws-sdk-go/aws/awsutil"
  "github.com/awslabs/aws-sdk-go/service/ec2"
  "github.com/awslabs/aws-sdk-go/service/iam"
)

var Commands = []cli.Command{
  commandList_ipaddress,
  commandEnter_instance,
  commandList_users,
}

var commandList_ipaddress = cli.Command{
  Name:  "list-ipaddress",
  Aliases:  []string{"li"},
  Usage: "get ipaddress list by NameTag",
  Description: `
`,
  Flags: []cli.Flag {
    cli.StringFlag{
      Name: "nametag, n",
      Usage: "NameTag",
    },
  },
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

var commandList_users = cli.Command{
  Name:  "list-users",
  Aliases:  []string{"lu"},
  Usage: "",
  Description: `
  List IAM Users
`,
  Action: doList_users,
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
  profile := c.GlobalString("profile")
  if profile == "" {
    fmt.Println("'--profile' is required")
    os.Exit(1)
  }
  name := c.String("nametag")
  if name == "" {
    fmt.Println("'--nametag' is required")
    os.Exit(1)
  }
  prov, _ := aws.ProfileCreds("", profile, 5 * time.Minute)
  svc := ec2.New(&aws.Config{Credentials: prov, Region: "ap-northeast-1"})
  params := ec2.DescribeInstancesInput{
    Filters: []ec2.Filter{
      Name: "tag:Name",
      Values: "*" + name + "*",
    },
  }
  res, err := svc.DescribeInstances(params)

  if awserr := aws.Error(err); awserr != nil {
    // A service error occurred.
    fmt.Println("Error:", awserr.Code, awserr.Message)
  } else if err != nil {
    // A non-service error occurred.
    panic(err)
  }

  for _, r := range res.Reservations {
    for _, i := range r.Instances {
      var nt string
      for _, t := range i.Tags {
        if *t.Key == "Name" {
          nt = *t.Value
          break
        }
      }
      fmt.Println(nt, *i.PrivateIPAddress)
    }
  }
  //fmt.Println(awsutil.StringValue(res))
}

func doList_users(c *cli.Context) {
  profile := c.GlobalString("profile")
  if profile == "" {
    fmt.Println("'--profile' is required")
    os.Exit(1)
  }
  prov, _ := aws.ProfileCreds("", profile, 5 * time.Minute)
  svc := iam.New(&aws.Config{Credentials: prov})
  params := &iam.ListUsersInput{
    //Marker:     aws.String("markerType"),
    //MaxItems:   aws.Long(1),
    //PathPrefix: aws.String("/"),
  }
  resp, err := svc.ListUsers(params)
  if awserr := aws.Error(err); awserr != nil {
    // A service error occurred.
    fmt.Println("Error:", awserr.Code, awserr.Message)
  } else if err != nil {
    // A non-service error occurred.
    panic(err)
  }
  fmt.Println(awsutil.StringValue(resp))
}

func doEnter_instance(c *cli.Context) {
}


