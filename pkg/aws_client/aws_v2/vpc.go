package aws_v2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	CON "github.com/openshift-qe/openshift-rosa-cli/pkg/constants"
	"github.com/openshift-qe/openshift-rosa-cli/pkg/log"
)

func (client *AwsV2Client) ListVPCByName(vpcName string) ([]types.Vpc, error) {
	vpcs := []types.Vpc{}
	filterKey := "tag:Name"
	filter := []types.Filter{
		types.Filter{
			Name:   &filterKey,
			Values: []string{vpcName},
		},
	}
	input := &ec2.DescribeVpcsInput{
		Filters: filter,
	}
	resp, err := client.ec2Client.DescribeVpcs(context.TODO(), input)
	if err != nil {
		return vpcs, err
	}
	vpcs = resp.Vpcs
	return vpcs, nil
}

func (client *AwsV2Client) CreateVpc(cidr string, name ...string) (*ec2.CreateVpcOutput, error) {
	vpcName := CON.VpcName
	if len(name) == 1 {
		vpcName = name[0]
	}
	tags := map[string]string{
		"Name":        vpcName,
		CON.QEFlagKey: CON.QEFLAG,
	}
	input := &ec2.CreateVpcInput{
		CidrBlock:         aws.String(cidr),
		DryRun:            nil,
		InstanceTenancy:   "",
		Ipv4IpamPoolId:    nil,
		Ipv4NetmaskLength: nil,
		TagSpecifications: nil,
	}

	resp, err := client.ec2Client.CreateVpc(context.TODO(), input)
	if err != nil {
		log.LogError("Create vpc error " + err.Error())
		return nil, err
	}
	log.LogInfo("Create vpc success " + *resp.Vpc.VpcId)
	err = client.WaitForResourceExisting(*resp.Vpc.VpcId, 10)
	if err != nil {
		return resp, err
	}

	client.TagResource(*resp.Vpc.VpcId, tags)

	log.LogInfo("Created vpc with ID " + *resp.Vpc.VpcId)
	return resp, err
}

// ModifyVpcDnsAttribute will modify the vpc attibutes
// dnsAttribute should be the value of "DnsHostnames" and "DnsSupport"
func (client *AwsV2Client) ModifyVpcDnsAttribute(vpcID string, dnsAttribute string, status bool) (*ec2.ModifyVpcAttributeOutput, error) {
	inputModifyVpc := &ec2.ModifyVpcAttributeInput{}

	if dnsAttribute == CON.VpcDnsHostnamesAttribute {
		inputModifyVpc = &ec2.ModifyVpcAttributeInput{
			VpcId:              aws.String(vpcID),
			EnableDnsHostnames: &types.AttributeBooleanValue{Value: aws.Bool(status)},
		}
	} else if dnsAttribute == CON.VpcDnsSupportAttribute {
		inputModifyVpc = &ec2.ModifyVpcAttributeInput{
			VpcId:            aws.String(vpcID),
			EnableDnsSupport: &types.AttributeBooleanValue{Value: aws.Bool(status)},
		}
	}

	resp, err := client.ec2Client.ModifyVpcAttribute(context.TODO(), inputModifyVpc)
	if err != nil {
		log.LogError("Modify vpc dns attribute failed " + err.Error())
		return nil, err
	}
	log.LogInfo("Modify vpc dns attribute success" + vpcID + dnsAttribute)
	return resp, err
}

func (client *AwsV2Client) DeleteVpc(vpcID string) (*ec2.DeleteVpcOutput, error) {
	input := &ec2.DeleteVpcInput{
		VpcId:  aws.String(vpcID),
		DryRun: nil,
	}

	resp, err := client.ec2Client.DeleteVpc(context.TODO(), input)
	if err != nil {
		log.LogError("Delete vpc %s failed "+err.Error(), vpcID)
		return nil, err
	}
	log.LogInfo("Delete vpc success " + vpcID)
	return resp, err

}
func (client *AwsV2Client) DescribeVPC(vpcID string) (types.Vpc, error) {
	var vpc types.Vpc
	input := &ec2.DescribeVpcsInput{
		VpcIds: []string{vpcID},
	}

	resp, err := client.ec2Client.DescribeVpcs(context.TODO(), input)
	if err != nil {
		return vpc, err
	}
	vpc = resp.Vpcs[0]
	return vpc, err
}
