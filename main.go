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

	flag.StringVar(&fCommand, "cmd", "", "the command to run :STOP, TERMINATE")
	flag.StringVar(&fName, "name", "*", "the desired machine tagged name")
	flag.Parse()

	checkEnv(env_region, env_key_id, env_key)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	ec2Cli := ec2.New(sess)

	fmt.Printf("Working with tagged name:%s\n", fName)
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

	var instances []*ec2.Instance
	resL := len(result.Reservations)
	fmt.Printf("Number of reservation :%d\n", resL)
	if resL > 0 {
		for i, r := range result.Reservations {
			fmt.Printf("Reservation index :%d\n", i)
			fmt.Printf("Located instances matching the tag :%d\n", len(r.Instances))
			for _, i := range r.Instances {
				instances = append(instances, i)
				var nt string
				for _, t := range i.Tags {
					if *t.Key == "Name" {
						nt = *t.Value
						break
					}
				}
				if i.PublicIpAddress != nil {
					fmt.Println(nt, *i.InstanceId, *i.State.Name, *i.PublicIpAddress)
				} else {
					fmt.Println(nt, *i.InstanceId, *i.State.Name)
				}
			}
		}
	} else {
		fmt.Println("No matching instance...")
	}

	if len(instances) > 0 && fCommand != "" {
		fmt.Printf("Running command :%s\n", fCommand)
		switch fCommand {
		case "STOP":
			if confirm(fCommand) {
				stop(ec2Cli, instances)
			}
		case "TERMINATE":
			if confirm(fCommand) {
				kill(ec2Cli, instances)
			}

		default:
			fmt.Printf("Unknown command :%s", fCommand)
		}
	}
}

func confirm(c string) bool {
	fmt.Printf("You sure you want to %s all the previously listed instances (Y/n)", c)
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
