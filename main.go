package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/aws/awserr"

	"gopkg.in/cheggaaa/pb.v1"
	"github.com/joho/godotenv"
	"log"

	"os"
	"time"
	"flag"
	"net/url"
	"strings"
)

const (
	Version = "v0.0.4"
	Author  = " by Valentyn Nastenko [versus.dev@gmail.com]"
)

func main() {

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
	for i := 0; i < count;  {
		//bar.Increment()
		bar.Add64(5)
		i = i + 5
		time.Sleep(1*time.Second)
	}
	bar.FinishPrint("The End!")
	}()

	log.Println(os.Getenv("AWS_ACCESS_KEY_ID"))

	if *flagDry == false {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})
	svc := ec2.New(sess)


		params := &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				&ec2.Filter{
					Name: aws.String("instance-state-name"),
					Values: []*string{
						aws.String("running"),
						aws.String("pending"),
					},
				},
			},
		}

		// TODO: Actually care if we can't connect to a host
		resp, err := svc.DescribeInstances(params)
		if err != nil {
		      panic(err)
		 }
		fmt.Println(resp)
		// Loop through the instances. They don't always have a name-tag so set it
		// to None if we can't find anything.
		for idx, _ := range resp.Reservations {
			for _, inst := range resp.Reservations[idx].Instances {

				// We need to see if the Name is one of the tags. It's not always
				// present and not required in Ec2.
				name := "None"
				for _, keys := range inst.Tags {
					if *keys.Key == "Name" {
						name = url.QueryEscape(*keys.Value)
					}
				}

				important_vals := []*string{
					inst.InstanceId,
					&name,
					inst.PrivateIpAddress,
					inst.InstanceType,
					inst.PublicIpAddress,

				}

				// Convert any nil value to a printable string in case it doesn't
				// doesn't exist, which is the case with certain values
				output_vals := []string{}
				for _, val := range important_vals {
					if val != nil {
						output_vals = append(output_vals, *val)
					} else {
						output_vals = append(output_vals, "None")
					}
				}
				// The values that we care about, in the order we want to print them
				fmt.Println(strings.Join(output_vals, " "))
			}
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
}

