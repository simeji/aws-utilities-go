package main

import (
	//"bufio"
	"bytes"
	"code.google.com/p/go.crypto/ssh"
	"fmt"
	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/aws/awsutil"
	"github.com/awslabs/aws-sdk-go/aws/credentials"
	"github.com/awslabs/aws-sdk-go/service/ec2"
	"github.com/awslabs/aws-sdk-go/service/iam"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"log"
	"os"
)

var Commands = []cli.Command{
	commandExec_instance,
	commandList_users,
	commandOperate_instance,
}

var commandExec_instance = cli.Command{
	Name:    "exec-instance",
	Aliases: []string{"e"},
	Usage:   "enter the instance and execute command",
	Description: `
`,
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "nametag, n", Usage: "[*required] NameTag"},
		&cli.StringFlag{Name: "command, c", Usage: "[*required] Execution commands at remote host"},
		&cli.BoolFlag{Name: "public, pub", Usage: "Then true, remote server is accessed by PublicIPAddress"},
		&cli.StringFlag{Name: "user, u", Usage: "UserName used for login to remote host"},
		&cli.StringFlag{Name: "id, i", Usage: "PrivateKeyFile used for login to remote host"},
		&cli.StringFlag{Name: "port, p", Usage: "server port used for login to remote host"},
	},
	Action: doExec_instance,
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

var commandOperate_instance = cli.Command{
	Name:    "operate-instance",
	Aliases: []string{"o", "op"},
	Usage:   "operate instance start/stop filtered by NameTag",
	Description: `
	aws-utilities -p {your_profile} op -n {NameTag} -m {start/stop/status}
`,
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "nametag, n", Usage: "NameTag [default: *]"},
		&cli.StringFlag{Name: "mode, m", Usage: "[*required] start/stop/status"},
	},
	Action: func(c *cli.Context) {

		name, err := getNametag(c)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		mode, merr := getMode(c)
		if merr != nil {
			fmt.Println(merr)
			os.Exit(1)
		}

		profile := getProfile(c)
		cred := getCredential(profile)
		svc := ec2.New(getConfig(cred))

		status := getModeStatus(mode)

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

		if status != "" {
			params.Filters = append(params.Filters,
				&ec2.Filter{
					Name: aws.String("instance-state-name"),
					Values: []*string{
						aws.String(status),
					},
				},
			)
		}

		res, err := svc.DescribeInstances(params)

		if awserr := aws.Error(err); awserr != nil {
			// A service error occurred.
			fmt.Println("Error:", awserr.Code, awserr.Message)
			os.Exit(1)
		} else if err != nil {
			// A non-service error occurred.
			panic(err)
			os.Exit(1)
		}

		var instanceIds []*string

		for _, r := range res.Reservations {
			for _, i := range r.Instances {
				var nt, privateip, publicip string
				for _, t := range i.Tags {
					if *t.Key == "Name" {
						nt = *t.Value
						break
					}
				}
				if i.PrivateIPAddress != nil {
					privateip = *i.PrivateIPAddress
				}
				if i.PublicIPAddress != nil {
					publicip = *i.PublicIPAddress
				}
				fmt.Println(
					nt,
					privateip,
					publicip,
					*i.InstanceID,
					*i.State.Name,
				)
				instanceIds = append(instanceIds, i.InstanceID)
			}
		}

		if len(instanceIds) < 1 {
			fmt.Println("No " + status + " instances..Search by [" + name + "]")
			os.Exit(0)
		}

		switch mode {
		case "start":
			input := &ec2.StartInstancesInput{
				InstanceIDs: instanceIds,
			}
			sOut, sErr := svc.StartInstances(input)
			if err := formatAwsError(sErr); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(awsutil.StringValue(sOut))
		case "stop":
			input := &ec2.StopInstancesInput{
				InstanceIDs: instanceIds,
			}
			sOut, sErr := svc.StopInstances(input)
			if err := formatAwsError(sErr); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(awsutil.StringValue(sOut))
		}
	},
}

func getProfile(c *cli.Context) (profile string) {
	profile = c.GlobalString("profile")
	if profile == "" {
		profile = "default"
	}
	return
}

func getCredential(profile string) *credentials.Credentials {
	return credentials.NewSharedCredentials("", profile)
}

func getMode(c *cli.Context) (mode string, err error) {
	err = nil
	modes := map[string]bool{
		"start":  true,
		"stop":   true,
		"status": true,
	}
	mode = c.String("mode")
	if mode == "" {
		err = fmt.Errorf("'--mode' is required")
	} else if !modes[mode] {
		err = fmt.Errorf("--mode must be included in [start|stop|status]")
	}
	return
}

func getModeStatus(mode string) string {
	modes := map[string]string{
		"start":  "stopped",
		"stop":   "running",
		"status": "",
	}
	return modes[mode]
}

func getNametag(c *cli.Context) (name string, err error) {
	err = nil
	name = c.String("nametag")
	if name == "" {
		err = fmt.Errorf("'--nametag' is required")
	}
	return
}

func getConfig(cred *credentials.Credentials) *aws.Config {
	return &aws.Config{Credentials: cred, Region: "ap-northeast-1"}
}

func formatAwsError(err error) error {
	var aerr error = nil
	if awserr := aws.Error(err); awserr != nil {
		// A service error occurred.
		aerr = fmt.Errorf("Error:", awserr.Code, awserr.Message)
	} else if err != nil {
		// A non-service error occurred.
		aerr = aerr
	}
	return aerr
}

func doList_users(c *cli.Context) {
	profile := getProfile(c)
	cred := getCredential(profile)
	svc := iam.New(getConfig(cred))
	params := &iam.ListUsersInput{}
	resp, err := svc.ListUsers(params)

	if awserr := aws.Error(err); awserr != nil {
		// A service error occurred.
		fmt.Println("Error:", awserr.Code, awserr.Message)
		os.Exit(1)
	} else if err != nil {
		// A non-service error occurred.
		panic(err)
	}
	fmt.Println(awsutil.StringValue(resp.Users))
}

func doExec_instance(c *cli.Context) {
	name, err := getNametag(c)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	command := c.String("command")
	if command == "" {
		fmt.Println("'--command' is required")
		os.Exit(1)
	}
	profile := getProfile(c)
	cred := getCredential(profile)
	svc := ec2.New(getConfig(cred))
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("tag:Name"),
				Values: []*string{
					aws.String(name),
				},
			},
			&ec2.Filter{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String("running"),
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

	var ip string
	r := res.Reservations[0]
	for _, i := range r.Instances {
		ip = *i.PrivateIPAddress
		if c.Bool("public") {
			ip = *i.PublicIPAddress
		}
		break
	}

	// default port
	port := "22"
	_p := c.String("port")
	if _p != "" {
		port = _p
	}
	// default pem
	idfile := os.Getenv("HOME") + "/.ssh/id_rsa"
	if c.String("id") != "" {
		idfile = c.String("id")
	}
	// defult user
	username := os.Getenv("USER")
	if c.String("user") != "" {
		username = c.String("user")
	}
	contents, err := ioutil.ReadFile(idfile)
	if err != nil {
		fmt.Println(contents, err)
		os.Exit(2)
	}
	signer, err := ssh.ParsePrivateKey(contents)
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	client, err := ssh.Dial("tcp", ip+":"+port, config)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}

	session, err := client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	// Request pseudo terminal
	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		log.Fatalf("request for pseudo terminal failed: %s", err)
	}

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(command); err != nil {
		panic("Failed to run: " + err.Error())
	}
	fmt.Println(b.String())
}
