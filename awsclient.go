package main

import (
	"log"
	"os"

	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type AWSClient struct {
	ec2Client *ec2.EC2
	Region    string
}

func NewAWSClient() (*AWSClient, error) {
	s := &AWSClient{}
	ses, _ := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})
	s.ec2Client = ec2.New(ses)
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

	resp, err := cl.ec2Client.DescribeTags(params)
	if err != nil {
		return nil, parseAwsError(err)
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
	log.Println("Adding tags for %v, as %v", resourceID, tags)
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

	_, err := cl.ec2Client.CreateTags(params)
	if err != nil {
		return parseAwsError(err)
	}
	return nil
}
