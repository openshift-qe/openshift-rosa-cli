package aws_v2

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"gitlab.cee.redhat.com/openshift-group-I/ocm_aws/pkg/log"
)

func (client *AwsV2Client) ListS3EndPointAssociation(vpcID string) ([]types.VpcEndpoint, error) {
	vpcFilterKey := "vpc-id"
	filters := []types.Filter{
		types.Filter{
			Name:   &vpcFilterKey,
			Values: []string{vpcID},
		},
	}

	input := ec2.DescribeVpcEndpointsInput{
		Filters: filters,
	}
	resp, err := client.ec2Client.DescribeVpcEndpoints(context.TODO(), &input)
	if err != nil {
		return nil, err
	}
	return resp.VpcEndpoints, err
}

func (client *AwsV2Client) DeleteVPCEndpoints(vpcID string) error {
	vpcEndpoints, err := client.ListS3EndPointAssociation(vpcID)
	if err != nil {
		return err
	}
	var endpoints = []string{}
	for _, ve := range vpcEndpoints {
		endpoints = append(endpoints, *ve.VpcEndpointId)
	}
	if len(endpoints) != 0 {
		input := &ec2.DeleteVpcEndpointsInput{
			VpcEndpointIds: endpoints,
		}
		_, err = client.ec2Client.DeleteVpcEndpoints(context.TODO(), input)
	}
	if err != nil {
		log.LogError("Delete vpc endpoints %s failed: %s", strings.Join(endpoints, ","), err.Error())
	} else {
		log.LogInfo("Delete vpc endpoints %s successfully", strings.Join(endpoints, ","))
	}
	return err
}
