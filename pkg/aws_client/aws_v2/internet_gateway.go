package aws_v2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/openshift-qe/openshift-rosa-cli/pkg/log"
)

func (client *AwsV2Client) CreateInternetGateway() (*ec2.CreateInternetGatewayOutput, error) {
	inputCreateInternetGateway := &ec2.CreateInternetGatewayInput{
		DryRun:            nil,
		TagSpecifications: nil,
	}
	respCreateInternetGateway, err := client.ec2Client.CreateInternetGateway(context.TODO(), inputCreateInternetGateway)
	if err != nil {
		log.LogError("Create igw error " + err.Error())
		return nil, err
	}
	log.LogInfo("Create igw success: " + *respCreateInternetGateway.InternetGateway.InternetGatewayId)
	return respCreateInternetGateway, err
}

func (client *AwsV2Client) AttachInternetGateway(internetGatewayID string, vpcID string) (*ec2.AttachInternetGatewayOutput, error) {

	input := &ec2.AttachInternetGatewayInput{
		InternetGatewayId: aws.String(internetGatewayID),
		VpcId:             aws.String(vpcID),
		DryRun:            nil,
	}
	resp, err := client.ec2Client.AttachInternetGateway(context.TODO(), input)
	if err != nil {
		log.LogError("Attach igw error " + err.Error())
		return nil, err
	}
	log.LogInfo("Attach igw success: " + internetGatewayID)
	return resp, err
}

func (client *AwsV2Client) DetachInternetGateway(internetGatewayID string, vpcID string) (*ec2.DetachInternetGatewayOutput, error) {
	input := &ec2.DetachInternetGatewayInput{
		InternetGatewayId: aws.String(internetGatewayID),
		VpcId:             aws.String(vpcID),
		DryRun:            nil,
	}
	resp, err := client.ec2Client.DetachInternetGateway(context.TODO(), input)
	if err != nil {
		log.LogError("Detach igw %s error  from vpc %s:"+err.Error(), internetGatewayID, vpcID)
		return nil, err
	}
	log.LogInfo("Detach igw %s success from vpc %s", internetGatewayID, vpcID)
	return resp, err
}
func (client *AwsV2Client) ListInternetGateWay(vpcID string) ([]types.InternetGateway, error) {
	vpcFilter := "attachment.vpc-id"
	filter := []types.Filter{
		types.Filter{
			Name: &vpcFilter,
			Values: []string{
				vpcID,
			},
		},
	}
	input := &ec2.DescribeInternetGatewaysInput{
		Filters: filter,
	}
	resp, err := client.ec2Client.DescribeInternetGateways(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return resp.InternetGateways, err
}
func (client *AwsV2Client) DeleteInternetGateway(internetGatewayID string) (*ec2.DeleteInternetGatewayOutput, error) {
	inputDeleteInternetGateway := &ec2.DeleteInternetGatewayInput{
		InternetGatewayId: aws.String(internetGatewayID),
		DryRun:            nil,
	}
	respDeleteInternetGateway, err := client.ec2Client.DeleteInternetGateway(context.TODO(), inputDeleteInternetGateway)
	if err != nil {
		log.LogError("Delete igw error " + err.Error())
		return nil, err
	}
	log.LogInfo("Delete igw success: " + internetGatewayID)
	return respDeleteInternetGateway, err
}
