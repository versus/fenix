package main

import (
	"log"

	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func GetTags(ec2Client *ec2.EC2, resourceID string) (map[string]string, error) {
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

	resp, err := ec2Client.DescribeTags(params)
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

func AddTags(ec2Client *ec2.EC2, resourceID string, tags map[string]string) error {
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

	_, err := ec2Client.CreateTags(params)
	if err != nil {
		return parseAwsError(err)
	}
	return nil
}
