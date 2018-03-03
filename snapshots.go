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
	"gopkg.in/cheggaaa/pb.v1"
)

type CreateSnapshotRequest struct {
	VolumeID    string
	Description string
	Tags        map[string]string
}

type SnapshotDescribeRequest struct {
	SnapshotId string
}

func CreateSnapshot(cl *AWSClient, request *CreateSnapshotRequest) (string, error) {
	params := &ec2.CreateSnapshotInput{
		VolumeId:    aws.String(request.VolumeID),
		Description: aws.String(request.Description),
	}
	resp, err := cl.ec2Client.CreateSnapshot(params)
	if err != nil {
		return "", parseAwsError(err)
	}
	if request.Tags != nil {
		if err := cl.AddTags(*resp.SnapshotId, request.Tags); err != nil {
			log.Println("Unable to tag %v with %v, but continue", *resp.SnapshotId, request.Tags)
		}
	}
	return *resp.SnapshotId, nil
}

func WaitForSnapshotComplete(ec2Client *ec2.EC2, snapshotID string) error {
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
		snapshots, err := ec2Client.DescribeSnapshots(snapInput)
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
