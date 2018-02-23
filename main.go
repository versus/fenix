package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"

	"log"

	"github.com/joho/godotenv"
	"gopkg.in/cheggaaa/pb.v1"

	"flag"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	Version = "v0.0.4"
	Author  = " by Valentyn Nastenko [versus.dev@gmail.com]"
)

//https://github.com/rancher/convoy/blob/master/ebs/ebs_service.go

func main() {
	var SnapshotId *string
	log.Println("fenix ", Version, Author)
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env. file")
	}

	flagInstanceId := flag.String("instance", "", "AWS instanceId to replicate")
	flagDry := flag.Bool("dry", false, "true is Dry Run")
	flag.Parse()

	log.Println(flagInstanceId)

	go func() {
		count := 100
		bar := pb.StartNew(count)
		bar.SetMaxWidth(80)
		bar.ShowTimeLeft = false
		for i := 0; i < count; {
			//bar.Increment()
			bar.Add64(5)
			i = i + 5
			time.Sleep(1 * time.Second)
		}
		bar.FinishPrint("The End!")
	}()

	log.Println(os.Getenv("AWS_ACCESS_KEY_ID"))

	if *flagDry == false {
		sess, err := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})
		svc := ec2.New(sess)

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

		resp, err := svc.DescribeInstances(params)
		if err != nil {
			panic(err)
		}
		fmt.Println(resp)
		// Loop through the instances. They don't always have a name-tag so set it
		// to None if we can't find anything.
		for idx := range resp.Reservations {
			for _, inst := range resp.Reservations[idx].Instances {

				// We need to see if the Name is one of the tags. It's not always
				// present and not required in Ec2.
				name := "None"
				for _, keys := range inst.Tags {
					if *keys.Key == "Name" {
						name = url.QueryEscape(*keys.Value)
					}
				}

				hdd := inst.BlockDeviceMappings
				log.Println("count block devices ", len(hdd))
				log.Println("dev 1: ", hdd[1].Ebs.VolumeId)
				input := &ec2.CreateSnapshotInput{
					Description: aws.String("This is data volume snapshot."),
					VolumeId:    hdd[1].Ebs.VolumeId,
				}

				//res, err := svc.CreateSnapshotRequest(input)
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
				SnapshotId = res.SnapshotId

				importantVals := []*string{
					inst.InstanceId,
					&name,
					inst.PrivateIpAddress,
					inst.InstanceType,
					inst.PublicIpAddress,
				}

				// Convert any nil value to a printable string in case it doesn't
				// doesn't exist, which is the case with certain values
				var outputVals []string
				for _, val := range importantVals {
					if val != nil {
						outputVals = append(outputVals, *val)
					} else {
						outputVals = append(outputVals, "None")
					}
				}
				// The values that we care about, in the order we want to print them
				fmt.Println(strings.Join(outputVals, " "))
			}
		}

		snapInput := &ec2.DescribeSnapshotsInput{
			SnapshotIds: []*string{
				//aws.String("snap-0f411b956abd7ece8"),
				SnapshotId,
			},
		}

		for i := 0; i < 10; i++ {
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
			if *snapshots.Snapshots[0].State == "completed" {
				break
			}
			time.Sleep(2 * time.Second)
		}

	}
}
