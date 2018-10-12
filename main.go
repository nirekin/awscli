package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {

	var fName string
	var fCommand string

	flag.StringVar(&fCommand, "cmd", "-", "the command to run : STOP, KILL")
	flag.StringVar(&fName, "name", "*", "the desired mchine tagged name")
	flag.Parse()

	checkEnv(env_region, env_key_id, env_key)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	ec2Cli := ec2.New(sess)

	fmt.Println("Working with tagged name:%s", fName)
	sTag := "tag:Name"
	inParam := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: &sTag,
				Values: []*string{
					&fName,
				},
			},
		},
	}

	result, err := ec2Cli.DescribeInstances(inParam)
	if err != nil {
		panic(err)
	}

	resL := len(result.Reservations)
	if resL > 0 {
		fmt.Printf("Located instances : %d\n", len(result.Reservations[0].Instances))
		for _, i := range result.Reservations[0].Instances {
			var nt string
			for _, t := range i.Tags {
				if *t.Key == "Name" {
					nt = *t.Value
					break
				}
			}
			fmt.Println(nt, *i.InstanceId, *i.State.Name)

		}
	} else {
		fmt.Println("No instance located")
	}

	if resL > 0 && len(result.Reservations[0].Instances) > 0 && fCommand != "" {
		fmt.Printf("Running command :%s\n", fCommand)
		switch fCommand {
		case "STOP":
			if confirm(fCommand) {
				stop(ec2Cli, result.Reservations[0].Instances)
			}
		case "KILL":
			if confirm(fCommand) {
				kill(ec2Cli, result.Reservations[0].Instances)
			}

		default:
			fmt.Printf("Unknown command  :%s", fCommand)
		}
	}
}

func confirm(c string) bool {
	fmt.Printf("Sure you want to %s all the previously listed instances (Y/n)", c)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text) == "Y"
}

func checkEnv(keys ...string) {
	for _, k := range keys {
		v := os.Getenv(k)
		if v == "" {
			fmt.Printf("Missing environment variable %s", k)
			os.Exit(1)
		}
	}
}
