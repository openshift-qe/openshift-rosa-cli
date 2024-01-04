package aws_v2

import (
	"context"

	// elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	// elbtypes "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	elbtypes "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing/types"
	"gitlab.cee.redhat.com/openshift-group-I/ocm_aws/pkg/log"
)

func (client *AwsV2Client) DescribeLoadBalancers(vpcID string) ([]elbtypes.LoadBalancerDescription, error) {

	listenedELB := []elbtypes.LoadBalancerDescription{}
	input := &elb.DescribeLoadBalancersInput{}
	resp, err := client.elbClient.DescribeLoadBalancers(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	// for _, lb := range resp.LoadBalancers {
	for _, lb := range resp.LoadBalancerDescriptions {

		// if *lb.VpcId == vpcID {
		if *lb.VPCId == vpcID {
			log.LogInfo("Got load balancer %s", *lb.LoadBalancerName)
			listenedELB = append(listenedELB, lb)
		}
	}

	return listenedELB, err
}

func (client *AwsV2Client) DeleteELB(ELB elbtypes.LoadBalancerDescription) error {
	log.LogInfo("Goint to delete ELB %s", *ELB.LoadBalancerName)

	deleteELBInput := &elb.DeleteLoadBalancerInput{
		// LoadBalancerArn: ELB.LoadBalancerArn,
		LoadBalancerName: ELB.LoadBalancerName,
	}
	_, err := client.elbClient.DeleteLoadBalancer(context.TODO(), deleteELBInput)
	return err
}
