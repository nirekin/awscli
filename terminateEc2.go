package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func kill(cli *ec2.EC2, is []*ec2.Instance) {
	for _, i := range is {
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
