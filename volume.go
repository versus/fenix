package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/versus/fenix/awscl"
	"github.com/versus/fenix/utils"
)

/*


type VolumeRequest struct {
	Size       int64
	IOPS       int64
	SnapshotID string
	VolumeType string
	Tags       map[string]string
	KmsKeyID   string
	Encrypted  bool
}

func (s *ebsService) waitForVolumeAttaching(volumeID string) error {
	var attachment *ec2.VolumeAttachment
	volume, err := s.GetVolume(volumeID)
	if err != nil {
		return err
	}
	for len(volume.Attachments) == 0 {
		log.Debugf("Retry to get attachment of volume")
		volume, err = s.GetVolume(volumeID)
		if err != nil {
			return err
		}
	}
	attachment = volume.Attachments[0]

	for *attachment.State == ec2.VolumeAttachmentStateAttaching {
		log.Debugf("Waiting for volume %v attaching", volumeID)
		volume, err := s.GetVolume(volumeID)
		if err != nil {
			return err
		}
		if len(volume.Attachments) != 0 {
			attachment = volume.Attachments[0]
		} else {
			return fmt.Errorf("Attaching failed for ", volumeID)
		}
	}
	if *attachment.State != ec2.VolumeAttachmentStateAttached {
		return fmt.Errorf("Cannot attach volume, final state %v", *attachment.State)
	}
	return nil
}

func (s *ebsService) CreateVolume(request *CreateEBSVolumeRequest) (string, error) {
	if request == nil {
		return "", fmt.Errorf("Invalid CreateEBSVolumeRequest")
	}
	size := request.Size
	iops := request.IOPS
	snapshotID := request.SnapshotID
	volumeType := request.VolumeType
	kmsKeyID := request.KmsKeyID

	// EBS size are in GB, we would round it up
	ebsSize := size / GB
	if size%GB > 0 {
		ebsSize += 1
	}

	params := &ec2.CreateVolumeInput{
		AvailabilityZone: aws.String(s.AvailabilityZone),
		Size:             aws.Int64(ebsSize),
		Encrypted:        aws.Bool(request.Encrypted),
	}

	if snapshotID != "" {
		params.SnapshotId = aws.String(snapshotID)
	} else if kmsKeyID != "" {
		params.KmsKeyId = aws.String(kmsKeyID)
		params.Encrypted = aws.Bool(true)
	}

	if volumeType != "" {
		if err := checkVolumeType(volumeType); err != nil {
			return "", err
		}
		if volumeType == "io1" && iops == 0 {
			return "", fmt.Errorf("Invalid IOPS for volume type io1")
		}
		if volumeType != "io1" && iops != 0 {
			return "", fmt.Errorf("IOPS only valid for volume type io1")
		}
		params.VolumeType = aws.String(volumeType)
		if iops != 0 {
			params.Iops = aws.Int64(iops)
		}
	}

	ec2Volume, err := s.ec2Client.CreateVolume(params)
	if err != nil {
		return "", parseAwsError(err)
	}

	volumeID := *ec2Volume.VolumeId
	if err = s.waitForVolumeTransition(volumeID, ec2.VolumeStateCreating, ec2.VolumeStateAvailable); err != nil {
		log.Debug("Failed to create volume: ", err)
		err = s.DeleteVolume(volumeID)
		if err != nil {
			log.Errorf("Failed deleting volume: %v", parseAwsError(err))
		}
		return "", fmt.Errorf("Failed creating volume with size %v and snapshot %v",
			size, snapshotID)
	}
	if request.Tags != nil {
		if err := s.AddTags(volumeID, request.Tags); err != nil {
			log.Warnf("Unable to tag %v with %v, but continue", volumeID, request.Tags)
		}
	}

	return volumeID, nil
}

*/

func GetVolume(cl *awscl.AWSClient, volumeID string) (*ec2.Volume, error) {
	params := &ec2.DescribeVolumesInput{
		VolumeIds: []*string{
			aws.String(volumeID),
		},
	}
	volumes, err := cl.Ec2Client.DescribeVolumes(params)
	if err != nil {
		return nil, utils.ParseAwsError(err)
	}
	if len(volumes.Volumes) != 1 {
		return nil, fmt.Errorf("Cannot find volume %v", volumeID)
	}
	return volumes.Volumes[0], nil
}
