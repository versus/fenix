package main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Ec2Instance struct {
	metadataClient *ec2metadata.EC2Metadata
	ec2Client      *ec2.EC2
	ses            *session.Session

	InstanceID       string
	Region           string
	AvailabilityZone string
}

func NewEc2Instance() (*Ec2Instance, error) {
	s := &Ec2Instance{}
	s.ses, _ = session.NewSession(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})
	s.ec2Client = ec2.New(s.ses)
	return s, nil
}
