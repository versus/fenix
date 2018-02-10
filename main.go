package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

func main() {

	/*
	aws_user := os.Getenv("AWS_ACCESS_KEY_ID")
	aws_pass := os.Getenv("AWS_SECRET_ACCESS_KEY")

	sess := aws.Creds(aws_user, aws_pass, "")
	svc := ec2.New(sess, "us-west-1", nil)
	*/

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	svc := ec2.New(sess)

	// Call to get detailed information on each instance
    result, err := svc.DescribeInstances(nil)
    if err != nil {
        fmt.Println("Error", err)
    } else {
        fmt.Println("Success", result)
    }

    /*
    input := &ec2.CreateSnapshotInput{
		Description: aws.String("This is data volume snapshot."),
		VolumeId:    aws.String("vol-05365510f241d4d94"),
	}

	res, err := svc.CreateSnapshot(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(res)
    */

	snapInput := &ec2.DescribeSnapshotsInput{
		SnapshotIds: []*string{
			aws.String("snap-0f411b956abd7ece8"),
		},
	}


	snapshots, err := svc.DescribeSnapshots(snapInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(snapshots)

	for indx, _ := range snapshots.Snapshots {
		snap := snapshots.Snapshots[indx]
		fmt.Println("snap SnapshotId =", snap.SnapshotId)
	}


}

