package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	"log"

	"github.com/joho/godotenv"

	"flag"
)

const (
	Version = "v0.0.5"
	Author  = " by Valentyn Nastenko [versus.dev@gmail.com]"
)

//https://github.com/rancher/convoy/blob/master/ebs/ebs_service.go

func main() {
	var SnapshotId string
	log.Println("fenix ", Version, Author)
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env. file")
	}

	flagInstanceId := flag.String("instance", "", "AWS instanceId to replicate")
	flag.Parse()

	log.Println(flagInstanceId)
	//log.Println(os.Getenv("AWS_ACCESS_KEY_ID"))

	awsClient, _ := NewEc2Instance()

	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String("running"),
					aws.String("pending"),
				},
			},
		},
	}

	resp, err := awsClient.ec2Client.DescribeInstances(params)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)

	for idx := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {

			hdd := inst.BlockDeviceMappings
			log.Println("count block devices ", len(hdd))
			log.Println("dev 1: ", hdd[1].Ebs.VolumeId)
			resurceId := *hdd[1].Ebs.VolumeId
			tags, err := GetTags(awsClient.ec2Client, resurceId)
			if err != nil {
				log.Println("Error get tags", err)
			}
			requestSnapshot := CreateSnapshotRequest{
				VolumeID:    resurceId,
				Description: "This is data volume snapshot.",
				Tags:        tags,
			}
			SnapshotId, err = CreateSnapshot(awsClient.ec2Client, &requestSnapshot)
			if err != nil {
				panic(err)
			}
			WaitForSnapshotComplete(awsClient.ec2Client, SnapshotId)
		}
	}

}
