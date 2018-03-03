package main

import (
	"log"

	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/versus/fenix/awscl"
	"github.com/versus/fenix/utils"
	"gopkg.in/cheggaaa/pb.v1"
)

type FenixSnapshot interface {
	CreateSnapshot(cl *awscl.AWSClient, request *CreateSnapshotRequest) (string, error)
	WaitForSnapshotComplete(cl *awscl.AWSClient, snapshotID string) error
}

type CreateSnapshotRequest struct {
	VolumeID    string
	Description string
	Tags        map[string]string
}

func CreateSnapshot(cl *awscl.AWSClient, request *CreateSnapshotRequest) (string, error) {
	params := &ec2.CreateSnapshotInput{
		VolumeId:    aws.String(request.VolumeID),
		Description: aws.String(request.Description),
	}
	resp, err := cl.Ec2Client.CreateSnapshot(params)
	if err != nil {
		return "", utils.ParseAwsError(err)
	}
	if request.Tags != nil {
		if err := cl.AddTags(*resp.SnapshotId, request.Tags); err != nil {
			log.Println("Unable to tag %v with %v, but continue", *resp.SnapshotId, request.Tags)
		}
	}
	return *resp.SnapshotId, nil
}

func WaitForSnapshotComplete(cl *awscl.AWSClient, snapshotID string) error {
	snapInput := &ec2.DescribeSnapshotsInput{
		SnapshotIds: []*string{
			//aws.String("snap-0f411b956abd7ece8"),
			&snapshotID,
		},
	}
	ticker := time.NewTicker(time.Second)
	count := 100
	bar := pb.StartNew(count)
	bar.SetMaxWidth(80)
	bar.ShowTimeLeft = false

	for _ = range ticker.C {
		snapshots, err := cl.Ec2Client.DescribeSnapshots(snapInput)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				fmt.Println(err.Error())
			}
			break
		}

		snapshot := snapshots.Snapshots[0]
		percent, err := strconv.ParseInt(strings.Replace(*snapshot.Progress, "%", "", -1), 10, 64)
		if err != nil {
			log.Fatal("Error get percent: ", err.Error())
		}
		bar.Set(int(percent))
		if *snapshot.State == "completed" {
			bar.Finish()
			break
		}
	}
	ticker.Stop()
	return nil
}
