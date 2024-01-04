package aws_v2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"gitlab.cee.redhat.com/openshift-group-I/ocm_aws/pkg/log"
)

func (client *AwsV2Client) CreateNatGateway(subnetID string, allocationID string, vpcID string) (*ec2.CreateNatGatewayOutput, error) {
	inputCreateNat := &ec2.CreateNatGatewayInput{
		SubnetId:          aws.String(subnetID),
		AllocationId:      aws.String(allocationID),
		ClientToken:       nil,
		ConnectivityType:  "",
		DryRun:            nil,
		TagSpecifications: nil,
	}
	respCreateNat, err := client.ec2Client.CreateNatGateway(context.TODO(), inputCreateNat)
	if err != nil {
		log.LogError("Create nat error " + err.Error())
		return nil, err
	}
	log.LogInfo("Create nat success: " + *respCreateNat.NatGateway.NatGatewayId)
	err = client.WaitForResourceExisting(*respCreateNat.NatGateway.NatGatewayId, 10*60)
	return respCreateNat, err
}

// DeleteNatGateway will wait for <timeout> seconds for nat gateway becomes status of deleted
func (client *AwsV2Client) DeleteNatGateway(natGatewayID string, timeout ...int) (*ec2.DeleteNatGatewayOutput, error) {
	inputDeleteNatGateway := &ec2.DeleteNatGatewayInput{
		NatGatewayId: aws.String(natGatewayID),
		DryRun:       nil,
	}
	respDeleteNatGateway, err := client.ec2Client.DeleteNatGateway(context.TODO(), inputDeleteNatGateway)
	if err != nil {
		log.LogError("Delete Nat Gateway error " + err.Error())
		return nil, err
	}
	timeoutTime := 60
	if len(timeout) != 0 {
		timeoutTime = timeout[0]
	}
	err = client.WaitForResourceDeleted(natGatewayID, timeoutTime)
	if err != nil {
		return respDeleteNatGateway, err
	}
	log.LogInfo("Delete Nat Gateway success " + *respDeleteNatGateway.NatGatewayId)
	return respDeleteNatGateway, err
}

func (client *AwsV2Client) ListNatGateWays(vpcID string) ([]types.NatGateway, error) {
	vpcFilter := "vpc-id"
	filter := []types.Filter{
		types.Filter{
			Name: &vpcFilter,
			Values: []string{
				vpcID,
			},
		},
	}
	input := &ec2.DescribeNatGatewaysInput{
		Filter: filter,
	}
	output, err := client.ec2Client.DescribeNatGateways(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return output.NatGateways, nil
}
