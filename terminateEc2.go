package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func kill(cli *ec2.EC2, is []*ec2.Instance) {
	for _, i := range is {
		if *i.State.Name == "terminated" {
			fmt.Printf("Instance :%s already terminated\n", *i.InstanceId)
		} else if *i.State.Name == "shutting-down" {
			fmt.Printf("Instance :%s already shutting-down\n", *i.InstanceId)
		} else {
			fmt.Printf("Killing :%s\n", *i.InstanceId)
			input := &ec2.TerminateInstancesInput{
				InstanceIds: []*string{
					i.InstanceId,
				},
				DryRun: aws.Bool(true),
			}
			result, err := cli.TerminateInstances(input)
			awsErr, ok := err.(awserr.Error)
			if ok && awsErr.Code() == "DryRunOperation" {
				input.DryRun = aws.Bool(false)
				result, err = cli.TerminateInstances(input)
				if err != nil {
					fmt.Println("Error", err)
				} else {
					fmt.Println("Success", result.String())
				}
			} else {
				fmt.Println("Error", err)
			}
		}
	}
}
