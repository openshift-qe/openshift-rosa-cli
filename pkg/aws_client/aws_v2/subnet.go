package aws_v2

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	CON "github.com/openshift-qe/openshift-rosa-cli/pkg/constants"
	"github.com/openshift-qe/openshift-rosa-cli/pkg/log"
)

func (client *AwsV2Client) CreateSubnet(vpcID string, region string, zone string, subnetCidr string) (*types.Subnet, error) {
	if region == "" {
		region = CON.DefaultAWSRegion
	}
	if zone == "" {
		zone = CON.DefaultAWSZone
	}

	input := &ec2.CreateSubnetInput{
		VpcId:              aws.String(vpcID),
		AvailabilityZone:   aws.String(fmt.Sprintf(region + zone)),
		AvailabilityZoneId: nil,
		CidrBlock:          aws.String(subnetCidr),
		DryRun:             nil,
		Ipv6CidrBlock:      nil,
		Ipv6Native:         nil,
		OutpostArn:         nil,
		TagSpecifications:  nil,
	}
	respCreateSubnet, err := client.ec2Client.CreateSubnet(context.TODO(), input)
	if err != nil {
		log.LogError("create subnet error " + err.Error())
		return nil, err
	}
	log.LogInfo("Created subnet %s for vpc %s", *respCreateSubnet.Subnet.SubnetId, vpcID)
	err = client.WaitForResourceExisting(*respCreateSubnet.Subnet.SubnetId, 4)
	if err != nil {
		return nil, err
	}
	return respCreateSubnet.Subnet, err
}

func (client *AwsV2Client) ListSubnetByVpcID(vpcID string) ([]types.Subnet, error) {
	subnetFilter := []types.Filter{
		{
			Name: aws.String("vpc-id"),
			Values: []string{
				vpcID,
			},
		},
	}

	return client.ListSubnetsByFilter(subnetFilter)
}

func (client *AwsV2Client) DeleteSubnet(subnetID string) (*ec2.DeleteSubnetOutput, error) {
	input := &ec2.DeleteSubnetInput{
		SubnetId: aws.String(subnetID),
		DryRun:   nil,
	}

	resp, err := client.ec2Client.DeleteSubnet(context.TODO(), input)
	if err != nil {
		log.LogError("Delete subnet error " + err.Error())
		return nil, err
	}
	log.LogInfo("Delete subnet success " + subnetID)
	return resp, err
}

func (client *AwsV2Client) ListSubnetDetail(subnetIDs ...string) ([]types.Subnet, error) {
	// subnetFilter := []types.Filter{types.Filter{Name: aws.String("vpc-id"), Values: []string{vpcID}}}
	var subs = []types.Subnet{}
	if len(subnetIDs) == 0 {
		return subs, nil
	}

	input := &ec2.DescribeSubnetsInput{
		DryRun:     nil,
		Filters:    nil,
		MaxResults: nil,
		NextToken:  nil,
		SubnetIds:  subnetIDs,
	}

	resp, err := client.ec2Client.DescribeSubnets(context.TODO(), input)

	if err != nil {
		return subs, err
	}
	subs = resp.Subnets
	return subs, nil
}

// List subnet by filters
func (client *AwsV2Client) ListSubnetsByFilter(filter []types.Filter) ([]types.Subnet, error) {
	input := &ec2.DescribeSubnetsInput{
		DryRun:     nil,
		Filters:    filter,
		MaxResults: nil,
		NextToken:  nil,
		SubnetIds:  nil,
	}

	resp, err := client.ec2Client.DescribeSubnets(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("describe subnet by filter error " + err.Error())
	}

	return resp.Subnets, err
}
