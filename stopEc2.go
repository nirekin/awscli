package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func stop(cli *ec2.EC2, is []*ec2.Instance) {

	for _, i := range is {
		if *i.State.Name == "stopped" {
			fmt.Printf("Instance :%s already stopped \n", *i.InstanceId)
		} else if *i.State.Name == "terminated" {
			fmt.Printf("Instance :%s already terminated\n", *i.InstanceId)
		} else if *i.State.Name == "shutting-down" {
			fmt.Printf("Instance :%s already shutting-down\n", *i.InstanceId)
		} else {
			fmt.Printf("Stopping :%s\n", *i.InstanceId)
			input := &ec2.StopInstancesInput{
				InstanceIds: []*string{
					i.InstanceId,
				},
				DryRun: aws.Bool(true),
			}
			result, err := cli.StopInstances(input)
			awsErr, ok := err.(awserr.Error)
			if ok && awsErr.Code() == "DryRunOperation" {
				input.DryRun = aws.Bool(false)
				result, err = cli.StopInstances(input)
				if err != nil {
					fmt.Println("Error", err)
				} else {
					fmt.Println("Success", result.StoppingInstances)
				}
			} else {
				fmt.Println("Error", err)
			}
		}

	}
}
