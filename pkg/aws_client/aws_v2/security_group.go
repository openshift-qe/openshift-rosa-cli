package aws_v2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	CON "github.com/openshift-qe/openshift-rosa-cli/pkg/constants"
	"github.com/openshift-qe/openshift-rosa-cli/pkg/log"
)

func (client *AwsV2Client) ListSecurityGroups(vpcID string) ([]types.SecurityGroup, error) {
	vpcFilter := "vpc-id"
	customizedSGs := []types.SecurityGroup{}
	filter := []types.Filter{
		types.Filter{
			Name: &vpcFilter,
			Values: []string{
				vpcID,
			},
		},
	}
	describeSGInput := &ec2.DescribeSecurityGroupsInput{
		Filters: filter,
	}
	output, err := client.ec2Client.DescribeSecurityGroups(context.TODO(), describeSGInput)
	if err != nil {
		return nil, err
	}
	for _, sg := range output.SecurityGroups {
		if *sg.GroupName == "default" && *sg.Description == "default VPC security group" {
			continue
		}
		customizedSGs = append(customizedSGs, sg)
	}
	return customizedSGs, nil
}

func (client *AwsV2Client) ReleaseInboundOutboundRules(sgID string) error {
	filterKey := "group-id"
	filter := []types.Filter{
		types.Filter{
			Name: &filterKey,
			Values: []string{
				sgID,
			},
		},
	}
	describeSGInput := &ec2.DescribeSecurityGroupRulesInput{
		Filters: filter,
	}
	resp, err := client.ec2Client.DescribeSecurityGroupRules(context.TODO(), describeSGInput)
	if err != nil {
		log.LogError("Describe  rules failed for SG %s: %s", sgID, err.Error())
		return err
	}
	rules := resp.SecurityGroupRules
	ingressRules := []string{}
	egressRules := []string{}
	for _, rule := range rules {
		if *rule.IsEgress {
			egressRules = append(egressRules, *rule.SecurityGroupRuleId)
			continue
		}
		ingressRules = append(ingressRules, *rule.SecurityGroupRuleId)

	}
	if len(ingressRules) != 0 {
		releaseIngressRuleInput := &ec2.RevokeSecurityGroupIngressInput{
			GroupId:              &sgID,
			SecurityGroupRuleIds: ingressRules,
		}
		_, err = client.ec2Client.RevokeSecurityGroupIngress(context.TODO(), releaseIngressRuleInput)
		if err != nil {
			log.LogError("Release inbound rules failed for SG %s: %s", sgID, err.Error())
			return err
		}
	}
	if len(egressRules) != 0 {
		releaseEgressRuleInput := &ec2.RevokeSecurityGroupEgressInput{
			GroupId:              &sgID,
			SecurityGroupRuleIds: egressRules,
		}
		_, err = client.ec2Client.RevokeSecurityGroupEgress(context.TODO(), releaseEgressRuleInput)
		if err != nil {
			log.LogError("Release outbound rules failed for SG %s: %s", sgID, err.Error())
			return err
		}
	}
	log.LogInfo("Release rules successfully for SG %s", sgID)
	return nil
}

func (client *AwsV2Client) DeleteSecurityGroup(groupID string) (*ec2.DeleteSecurityGroupOutput, error) {

	err := client.ReleaseInboundOutboundRules(groupID)
	if err != nil {
		return nil, err
	}

	input := &ec2.DeleteSecurityGroupInput{
		DryRun:    nil,
		GroupId:   aws.String(groupID),
		GroupName: nil,
	}

	resp, err := client.ec2Client.DeleteSecurityGroup(context.TODO(), input)
	if err != nil {
		log.LogError("Delete security group %s failed %s", groupID, err.Error())
		return nil, err
	}
	log.LogInfo("Delete security group %s success ", groupID)
	return resp, err
}
func (client *AwsV2Client) AuthorizeSecurityGroupIngress(groupID string, cidr string, protocol string, fromPort int32, toPort int32) (*ec2.AuthorizeSecurityGroupIngressOutput, error) {
	input := &ec2.AuthorizeSecurityGroupIngressInput{
		CidrIp:                     aws.String(cidr),
		DryRun:                     nil,
		FromPort:                   aws.Int32(fromPort),
		GroupId:                    aws.String(groupID),
		GroupName:                  nil,
		IpPermissions:              nil,
		IpProtocol:                 aws.String(protocol),
		SourceSecurityGroupName:    nil,
		SourceSecurityGroupOwnerId: nil,
		TagSpecifications:          nil,
		ToPort:                     aws.Int32(toPort),
	}

	resp, err := client.ec2Client.AuthorizeSecurityGroupIngress(context.TODO(), input)
	if err != nil {
		log.LogError("Authorize security group failed " + err.Error())
		return nil, err
	}
	log.LogInfo("Authorize security group success " + groupID)
	return resp, err
}

func (client *AwsV2Client) CreateSecurityGroup(vpcID string, groupName string, sgDescription string) (*ec2.CreateSecurityGroupOutput, error) {
	input := &ec2.CreateSecurityGroupInput{
		Description:       aws.String(sgDescription),
		GroupName:         aws.String(groupName),
		DryRun:            nil,
		TagSpecifications: nil,
		VpcId:             aws.String(vpcID),
	}

	resp, err := client.ec2Client.CreateSecurityGroup(context.TODO(), input)
	if err != nil {
		log.LogError("Create security group failed " + err.Error())
		return nil, err
	}
	log.LogInfo("Create security group %s success for %s", *resp.GroupId, vpcID)
	err = client.WaitForResourceExisting(*resp.GroupId, 4)
	tags := map[string]string{
		"Name": CON.AdditionalSecurityGroupName,
	}
	_, err = client.TagResource(*resp.GroupId, tags)
	if err != nil {
		log.LogError("Created tagged failed %s", err)
	}
	log.LogInfo("Created tagged security group with ID %s", *resp.GroupId)
	return resp, err
}
