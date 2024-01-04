package aws_v2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/openshift-qe/openshift-rosa-cli/pkg/log"
)

func (client *AwsV2Client) ListNetWorkAcls(vpcID string) ([]types.NetworkAcl, error) {
	vpcFilter := "vpc-id"
	customizedAcls := []types.NetworkAcl{}
	filter := []types.Filter{
		types.Filter{
			Name: &vpcFilter,
			Values: []string{
				vpcID,
			},
		},
	}
	describeACLInput := &ec2.DescribeNetworkAclsInput{
		Filters: filter,
	}
	output, err := client.ec2Client.DescribeNetworkAcls(context.TODO(), describeACLInput)
	if err != nil {
		return nil, err
	}
	for _, acl := range output.NetworkAcls {
		customizedAcls = append(customizedAcls, acl)
	}
	return customizedAcls, nil
}

// RuleAction : deny/allow
// Protocol: TCP --> 6
func (client *AwsV2Client) AddNetworkAclEntry(networkAclId string, egress bool, protocol string, ruleAction string, ruleNumber int32, fromPort int32, toPort int32, cidrBlock string) (*ec2.CreateNetworkAclEntryOutput, error) {
	input := &ec2.CreateNetworkAclEntryInput{
		Egress:       aws.Bool(egress),
		NetworkAclId: aws.String(networkAclId),
		Protocol:     aws.String(protocol),
		RuleAction:   types.RuleAction(ruleAction),
		RuleNumber:   aws.Int32(ruleNumber),
		CidrBlock:    aws.String(cidrBlock),
		PortRange: &types.PortRange{
			From: aws.Int32(fromPort),
			To:   aws.Int32(toPort),
		},
	}
	resp, err := client.ec2Client.CreateNetworkAclEntry(context.TODO(), input)
	if err != nil {
		log.LogError("Create NetworkAcl rule failed " + err.Error())
		return nil, err
	}
	log.LogInfo("Create NetworkAcl rule success " + networkAclId)
	return resp, err
}

func (client *AwsV2Client) DeleteNetworkAclEntry(networkAclId string, egress bool, ruleNumber int32) (*ec2.DeleteNetworkAclEntryOutput, error) {
	input := &ec2.DeleteNetworkAclEntryInput{
		Egress:       aws.Bool(egress),
		NetworkAclId: aws.String(networkAclId),
		RuleNumber:   aws.Int32(ruleNumber),
	}
	resp, err := client.ec2Client.DeleteNetworkAclEntry(context.TODO(), input)
	if err != nil {
		log.LogError("Delete NetworkAcl rule failed " + err.Error())
		return nil, err
	}
	log.LogInfo("Delete NetworkAcl rule success " + networkAclId)
	return resp, err

}
