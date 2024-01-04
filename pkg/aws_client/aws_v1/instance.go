package aws_v1

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
)

func GetInstanceName(instance *ec2.Instance) string {
	tags := instance.Tags
	for _, tag := range tags {
		if *tag.Key == "Name" {
			return *tag.Value
		}
	}
	return ""
}

// GetInstancesByInfraID will return the instances with tag tag:kubernetes.io/cluster/<infraID>
func (client *AWSClient) GetInstancesByInfraID(infraID string) ([]*ec2.Instance, error) {
	filter := &ec2.Filter{
		Name: aws.String("tag:kubernetes.io/cluster/" + infraID),
		Values: []*string{
			aws.String("owned"),
		},
	}
	output, err := client.ec2Client.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			filter,
		},
		MaxResults: aws.Int64(100),
	})
	if err != nil {
		return nil, err
	}
	var instances []*ec2.Instance
	for _, reservation := range output.Reservations {
		instances = append(instances, reservation.Instances...)
	}
	return instances, err
}

func (client *AWSClient) GetInstanceRoleName(instance *ec2.Instance) (roleName string, err error) {
	instanceProfileName := strings.Split(*instance.IamInstanceProfile.Arn, "/")[1]
	getter := iam.GetInstanceProfileInput{
		InstanceProfileName: &instanceProfileName,
	}
	out, err := client.iamClient.GetInstanceProfile(&getter)
	if err != nil {
		return
	}
	roleName = *out.InstanceProfile.Roles[0].RoleName
	return
}
func (client *AWSClient) ListAvaliableRegionsFromAWS() ([]*ec2.Region, error) {
	optInStatus := "opt-in-status"
	optInNotRequired := "opt-in-not-required"
	optIn := "opted-in"
	filter := ec2.Filter{Name: &optInStatus, Values: []*string{&optInNotRequired, &optIn}}

	output, err := client.ec2Client.DescribeRegions(&ec2.DescribeRegionsInput{
		Filters: []*ec2.Filter{
			&filter,
		},
	})
	if err != nil {
		return nil, err
	}
	return output.Regions, err
}

func (client *AWSClient) ListMachineTypesPerRegion(regionID string) ([]string, error) {
	filter := ec2.Filter{
		Name:   aws.String("location"),
		Values: []*string{aws.String(regionID)},
	}
	params := &ec2.DescribeInstanceTypeOfferingsInput{
		Filters: []*ec2.Filter{
			&filter,
		},
	}
	out, err := client.ec2Client.DescribeInstanceTypeOfferings(params)
	if err != nil {
		return nil, err
	}
	fmt.Println(out.InstanceTypeOfferings)
	machineTypeList := make([]string, len(out.InstanceTypeOfferings))
	for i, v := range out.InstanceTypeOfferings {
		machineTypeList[i] = *v.InstanceType
	}
	return machineTypeList, nil

}
