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
	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-west-1")})
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
}

