package aws_v2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"gitlab.cee.redhat.com/openshift-group-I/ocm_aws/pkg/log"
)

func (awsClient AwsV2Client) CreateHostedZone(hostedZoneName string, vpcID string, private bool) (*route53.CreateHostedZoneOutput, error) {
	input := &route53.CreateHostedZoneInput{
		Name: &hostedZoneName,
		HostedZoneConfig: &types.HostedZoneConfig{
			PrivateZone: private,
		},
	}
	if vpcID != "" {
		vpc := &types.VPC{
			VPCId: &vpcID,
		}
		input.VPC = vpc
	}
	resp, err := awsClient.route53Client.CreateHostedZone(context.TODO(), input)
	if err != nil {
		log.LogError("Create hosted zone failed for vpc %s with name %s: ", vpcID, hostedZoneName, err.Error())
	} else {
		log.LogError("Create hosted zone succeed for vpc %s with name %s: ", vpcID, hostedZoneName, err.Error())
	}
	return resp, err
}

func (awsClient AwsV2Client) GetHostedZone(hostedZoneID string) (*route53.GetHostedZoneOutput, error) {
	input := &route53.GetHostedZoneInput{
		Id: &hostedZoneID,
	}

	return awsClient.route53Client.GetHostedZone(context.TODO(), input)
}

func (awsClient AwsV2Client) ListHostedZoneByDNSName(hostedZoneName string) (*route53.ListHostedZonesByNameOutput, error) {
	var maxItems int32 = 1
	input := &route53.ListHostedZonesByNameInput{
		DNSName:  &hostedZoneName,
		MaxItems: &maxItems,
	}

	return awsClient.route53Client.ListHostedZonesByName(context.TODO(), input)
}
