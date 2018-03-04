package awscl

import (
	"log"
	"os"

	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/versus/fenix/utils"
)

type AWSClient struct {
	Ec2Client    *ec2.EC2
	Region       string
	SourceVolume *ec2.Volume
	TargetVolume *ec2.Volume
	Snapshot     *ec2.Snapshot
}

func NewAWSClient() (*AWSClient, error) {
	s := &AWSClient{}
	ses, _ := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})
	s.Ec2Client = ec2.New(ses)
	return s, nil
}

func (cl *AWSClient) GetTags(resourceID string) (map[string]string, error) {
	params := &ec2.DescribeTagsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("resource-id"),
				Values: []*string{
					aws.String(resourceID),
				},
			},
		},
	}

	resp, err := cl.Ec2Client.DescribeTags(params)
	if err != nil {
		return nil, utils.ParseAwsError(err)
	}

	result := map[string]string{}
	if resp.Tags == nil {
		return result, nil
	}

	for _, tag := range resp.Tags {
		if *tag.ResourceId != resourceID {
			return nil, fmt.Errorf("BUG: why the result is not related to what I asked for?")
		}
		result[*tag.Key] = *tag.Value
	}
	return result, nil
}

func (cl *AWSClient) AddTags(resourceID string, tags map[string]string) error {
	if tags == nil {
		return nil
	}
	log.Println("Adding tags ", resourceID, tags)
	params := &ec2.CreateTagsInput{
		Resources: []*string{
			aws.String(resourceID),
		},
	}
	ec2Tags := []*ec2.Tag{}
	for k, v := range tags {
		tag := &ec2.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}
		ec2Tags = append(ec2Tags, tag)
	}
	params.Tags = ec2Tags

	_, err := cl.Ec2Client.CreateTags(params)
	if err != nil {
		return utils.ParseAwsError(err)
	}
	return nil
}
