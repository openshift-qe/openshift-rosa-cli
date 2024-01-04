package aws_v2

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func (client *AwsV2Client) DescribeVolumeByID(volumeID string) (*ec2.DescribeVolumesOutput, error) {

	output, err := client.ec2Client.DescribeVolumes(context.TODO(), &ec2.DescribeVolumesInput{
		VolumeIds: []string{volumeID},
	})

	if err != nil {
		fmt.Println("Got error describe volume: ", err)
	}
	return output, err
}
