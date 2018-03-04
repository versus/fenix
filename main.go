package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	"log"

	"github.com/joho/godotenv"

	"flag"

	"github.com/versus/fenix/awscl"
)

const (
	Version = "v0.0.5"
	Author  = " by Valentyn Nastenko [versus.dev@gmail.com]"
)

//https://github.com/rancher/convoy/blob/master/ebs/ebs_service.go

func main() {
	log.Println("fenix ", Version, Author)
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env. file")
	}

	flagInstanceId := flag.String("instance", "", "AWS instanceId for replicate")
	flag.Parse()

	log.Println(flagInstanceId)
	//log.Println(os.Getenv("AWS_ACCESS_KEY_ID"))

	awsClient, err := awscl.NewAWSClient()
	if err != nil {
		log.Fatal("FATAL: error create connect to aws Error:", err)
	}

	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String("running"),
					aws.String("pending"),
					aws.String("stopped"),
				},
			},
		},
	}

	resp, err := awsClient.Ec2Client.DescribeInstances(params)
	if err != nil {
		panic(err)
	}

	for idx := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {

			hdd := inst.BlockDeviceMappings
			log.Println("count block devices ", len(hdd))
			log.Println("dev 1: ", hdd[1].Ebs)
			resurceId := *hdd[1].Ebs.VolumeId

			awsClient.SourceVolume, err = GetVolume(awsClient, resurceId)
			if err != nil {
				log.Println("Error get volume ", err)
			}

			log.Println("Volume: ", awsClient.SourceVolume)

			tags, err := awsClient.GetTags(resurceId)
			if err != nil {
				log.Println("Error get tags", err)
			}
			requestSnapshot := CreateSnapshotRequest{
				VolumeID:    resurceId,
				Description: "This is data volume snapshot.",
				Tags:        tags,
			}

			awsClient.Snapshot, err = CreateSnapshot(awsClient, &requestSnapshot)
			if err != nil {
				panic(err)
			}
			WaitForSnapshotComplete(awsClient)
		}
	}
}
