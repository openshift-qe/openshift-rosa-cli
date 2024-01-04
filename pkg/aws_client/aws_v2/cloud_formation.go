package aws_v2

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"gitlab.cee.redhat.com/openshift-group-I/ocm_aws/pkg/log"
)

// Describe cloudformation by name
func (client *AwsV2Client) DescibeCloudFormationByName(cfName string) (*cloudformation.DescribeStackResourcesOutput, error) {
	input := &cloudformation.DescribeStackResourcesInput{
		StackName: &cfName,
	}
	resp, err := client.stackFormationClient.DescribeStackResources(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("describe cloudformation by filter error %s", err.Error())
	}
	log.LogInfo("%s", resp.ResultMetadata.Get("LogicalResourceId"))
	return resp, err
}

// List cloudformation stack resurouce by name
func (client *AwsV2Client) ListCloudFormationStackResourceByName(cfName string) (*cloudformation.ListStackResourcesOutput, error) {
	input := &cloudformation.ListStackResourcesInput{
		StackName: &cfName,
	}
	resp, err := client.stackFormationClient.ListStackResources(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("list cloudformation by filter error %s", err.Error())
	}
	log.LogInfo("%s", resp.ResultMetadata.Get("LogicalResourceId"))
	return resp, err
}
