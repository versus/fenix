package utils

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
)

func ParseAwsError(err error) error {
	if err == nil {
		return nil
	}
	if awsErr, ok := err.(awserr.Error); ok {
		message := fmt.Sprintln("AWS Error: ", awsErr.Code(), awsErr.Message(), awsErr.OrigErr())
		if reqErr, ok := err.(awserr.RequestFailure); ok {
			message += fmt.Sprintln(reqErr.StatusCode(), reqErr.RequestID())
		}
		return fmt.Errorf(message)
	}
	return err
}
