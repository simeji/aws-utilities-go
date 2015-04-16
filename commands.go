package main

import (
	"fmt"
	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/aws/awsutil"
	"github.com/awslabs/aws-sdk-go/service/ec2"
	"github.com/awslabs/aws-sdk-go/service/iam"
	"github.com/codegangsta/cli"
	"log"
	"os"
	"time"
)

var Commands = []cli.Command{
	commandList_ipaddress,
	commandEnter_instance,
	commandList_users,
}

var commandList_ipaddress = cli.Command{
	Name:    "list-ipaddress",
	Aliases: []string{"l"},
	Usage:   "list ipaddress filtered by NameTag",
	Description: `
`,
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "nametag, n", Usage: "NameTag"},
		&cli.BoolFlag{Name: "all, a", Usage: "Get all status instances"},
	},
	Action: doList_ipaddress,
}

var commandEnter_instance = cli.Command{
	Name:    "enter-instance",
	Aliases: []string{"e"},
	Usage:   "enter the instance filtered by NameTag",
	Description: `
	aaaa hogehoge mogemoge
`,
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "nametag, n", Usage: "NameTag"},
	},
	Action: doEnter_instance,
}

var commandList_users = cli.Command{
	Name:    "list-users",
	Aliases: []string{"u"},
	Usage:   "list all iam users",
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
	prov, _ := aws.ProfileCreds("", profile, 5*time.Minute)
	svc := ec2.New(&aws.Config{Credentials: prov, Region: "ap-northeast-1"})
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("tag:Name"),
				Values: []*string{
					aws.String(name),
				},
			},
		},
	}
	if c.Bool("all") == false {
		sf := &ec2.Filter{
			Name: aws.String("instance-state-name"),
			Values: []*string{
				aws.String("running"),
			},
		}
		params.Filters = append(params.Filters, sf)
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
			fmt.Println(nt, *i.PrivateIPAddress, *i.State.Name)
		}
	}
}

func doList_users(c *cli.Context) {
	profile := c.GlobalString("profile")
	if profile == "" {
		fmt.Println("'--profile' is required")
		os.Exit(1)
	}
	prov, _ := aws.ProfileCreds("", profile, 5*time.Minute)
	svc := iam.New(&aws.Config{Credentials: prov})
	params := &iam.ListUsersInput{
	//Marker:		 aws.String("markerType"),
	//MaxItems:	 aws.Long(1),
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
	fmt.Println(awsutil.StringValue(resp.Users))
}

func doEnter_instance(c *cli.Context) {
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
	prov, _ := aws.ProfileCreds("", profile, 5*time.Minute)
	svc := ec2.New(&aws.Config{Credentials: prov, Region: "ap-northeast-1"})
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("tag:Name"),
				Values: []*string{
					aws.String(name),
				},
			},
		},
	}
	res, err := svc.DescribeInstances(params)

	if awserr := aws.Error(err); awserr != nil {
		fmt.Println("Error:", awserr.Code, awserr.Message)
	} else if err != nil {
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
			fmt.Println(nt, *i.PrivateIPAddress, *i.State.Name)
		}
	}
}
