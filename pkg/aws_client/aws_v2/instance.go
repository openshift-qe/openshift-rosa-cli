package aws_v2

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/openshift-qe/openshift-rosa-cli/pkg/log"
)

func (client *AwsV2Client) LaunchInstance(subnetID string, imageID string, count int, instanceType string, keyName string, securityGroupIds []string, wait bool) (*ec2.RunInstancesOutput, error) {
	input := &ec2.RunInstancesInput{
		ImageId:          aws.String(imageID),
		MinCount:         aws.Int32(int32(count)),
		MaxCount:         aws.Int32(int32(count)),
		InstanceType:     types.InstanceType(instanceType),
		KeyName:          aws.String(keyName),
		SecurityGroupIds: securityGroupIds,
		SubnetId:         &subnetID,
	}
	output, err := client.ec2Client.RunInstances(context.TODO(), input)
	if wait && err == nil {
		instanceIDs := []string{}
		for _, instance := range output.Instances {
			instanceIDs = append(instanceIDs, *instance.InstanceId)
		}
		log.LogInfo("Waiting for below instances ready: %s", strings.Join(instanceIDs, "，"))
		_, err = client.WaitForInstancesRunning(instanceIDs, 10)
		if err != nil {
			log.LogError("Error happened for instance running: %s", err)
		} else {
			log.LogInfo("All instances running")
		}
	}
	return output, err
}

// ListInstance pass parameter like
// map[string][]string{"vpc-id":[]string{"<id>" }}, map[string][]string{"tag:Name":[]string{"<value>" }}
func (client *AwsV2Client) ListInstance(filters ...map[string][]string) ([]types.Instance, error) {
	FilterInput := []types.Filter{}
	for _, filter := range filters {
		for k, v := range filter {
			awsFilter := types.Filter{
				Name:   &k,
				Values: v,
			}
			FilterInput = append(FilterInput, awsFilter)
		}
	}
	getInstanceInput := &ec2.DescribeInstancesInput{
		Filters: FilterInput,
	}
	resp, err := client.EC2().DescribeInstances(context.TODO(), getInstanceInput)
	if err != nil {
		log.LogError("List instances failed with filters %v: %s", filters, err)
	}
	var instances []types.Instance
	for _, reserv := range resp.Reservations {
		instances = append(instances, reserv.Instances...)
	}
	return instances, err
}

func (client *AwsV2Client) WaitForInstanceReady(instanceID string, timeout time.Duration) error {
	instanceIDs := []string{
		instanceID,
	}
	log.LogInfo("Waiting for below instances ready: %s ", strings.Join(instanceIDs, "|"))
	_, err := client.WaitForInstancesRunning(instanceIDs, 10)
	return err
}

func (client *AwsV2Client) CheckInstanceState(instanceIDs ...string) (*ec2.DescribeInstanceStatusOutput, error) {
	log.LogInfo("Check instances status of %s", strings.Join(instanceIDs, ","))
	includeAll := true
	input := &ec2.DescribeInstanceStatusInput{
		InstanceIds:         instanceIDs,
		IncludeAllInstances: &includeAll,
	}
	output, err := client.ec2Client.DescribeInstanceStatus(context.TODO(), input)
	return output, err
}

// timeout indicates the minutes
func (client *AwsV2Client) WaitForInstancesRunning(instanceIDs []string, timeout time.Duration) (allRunning bool, err error) {
	startTime := time.Now()

	// output, err := client.CheckInstanceState(instanceIDs...)
	// if err != nil {
	// 	return
	// }
	for time.Now().Before(startTime.Add(timeout * time.Minute)) {
		allRunning = true
		output, err := client.CheckInstanceState(instanceIDs...)
		fmt.Println(output.InstanceStatuses)
		if err != nil {
			log.LogError("Error happened when describe instant status: %s", strings.Join(instanceIDs, ","))
			return false, err
		}
		if len(output.InstanceStatuses) == 0 {
			log.LogWarning("Instance status description for %s is 0", strings.Join(instanceIDs, ","))
		}
		for _, ins := range output.InstanceStatuses {
			log.LogInfo("Instance ID %s is in status of %s", *ins.InstanceId, ins.InstanceStatus.Status)
			log.LogInfo("Instance ID %s is in state of %s", *ins.InstanceId, ins.InstanceState.Name)
			if ins.InstanceState.Name != types.InstanceStateNameRunning && ins.InstanceStatus.Status != types.SummaryStatusOk {
				allRunning = false
			}

		}
		if allRunning {
			return true, nil
		}
		time.Sleep(time.Minute)
	}
	err = fmt.Errorf("timeout for waiting instances running")
	return
}
func (client *AwsV2Client) WaitForInstancesTerminated(instanceIDs []string, timeout time.Duration) (allTerminated bool, err error) {
	startTime := time.Now()
	for time.Now().Before(startTime.Add(timeout * time.Minute)) {
		allTerminated = true
		output, err := client.CheckInstanceState(instanceIDs...)
		if err != nil {
			log.LogError("Error happened when describe instant status: %s", strings.Join(instanceIDs, ","))
			return false, err
		}
		if len(output.InstanceStatuses) == 0 {
			log.LogWarning("Instance status description for %s is 0", strings.Join(instanceIDs, ","))
		}
		for _, ins := range output.InstanceStatuses {
			log.LogInfo("Instance ID %s is in status of %s", *ins.InstanceId, ins.InstanceStatus.Status)
			log.LogInfo("Instance ID %s is in state of %s", *ins.InstanceId, ins.InstanceState.Name)
			if ins.InstanceState.Name != types.InstanceStateNameTerminated {
				allTerminated = false
			}

		}
		if allTerminated {
			return true, nil
		}
		time.Sleep(time.Minute)
	}
	err = fmt.Errorf("timeout for waiting instances terminated")
	return

}

// Search instance types for specified region/availability zones
func (client *AwsV2Client) ListAvaliableInstanceTypesForRegion(region string, availabilityZones ...string) ([]string, error) {
	var params *ec2.DescribeInstanceTypeOfferingsInput
	if len(availabilityZones) > 0 {
		params = &ec2.DescribeInstanceTypeOfferingsInput{
			Filters:      []types.Filter{{Name: aws.String("location"), Values: availabilityZones}},
			LocationType: types.LocationTypeAvailabilityZone,
		}
	} else {
		params = &ec2.DescribeInstanceTypeOfferingsInput{
			Filters: []types.Filter{{Name: aws.String("location"), Values: []string{region}}},
		}
	}
	var instanceTypes []types.InstanceTypeOffering
	paginator := ec2.NewDescribeInstanceTypeOfferingsPaginator(client.ec2Client, params)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}
		instanceTypes = append(instanceTypes, page.InstanceTypeOfferings...)
	}
	machineTypeList := make([]string, len(instanceTypes))
	for i, v := range instanceTypes {
		machineTypeList[i] = string(v.InstanceType)
	}
	return machineTypeList, nil
}

// List avaliablezone for specific region
// zone type are: local-zone/availability-zone/wavelength-zone
func (client *AwsV2Client) ListAvaliableZonesForRegion(region string, zoneType string) ([]string, error) {
	var zones []string
	availabilityZones, err := client.ec2Client.DescribeAvailabilityZones(context.TODO(), &ec2.DescribeAvailabilityZonesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("region-name"),
				Values: []string{region},
			},
			{
				Name:   aws.String("zone-type"),
				Values: []string{zoneType},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	if len(availabilityZones.AvailabilityZones) < 1 {
		return zones, nil
	}

	for _, v := range availabilityZones.AvailabilityZones {
		zones = append(zones, *v.ZoneName)
	}
	return zones, nil
}
func (client *AwsV2Client) TerminateInstances(instanceIDs []string, wait bool, timeout time.Duration) error {
	if len(instanceIDs) == 0 {
		log.LogInfo("Got no instances to terminate.")
		return nil
	}
	terminateInput := &ec2.TerminateInstancesInput{
		InstanceIds: instanceIDs,
	}
	_, err := client.EC2().TerminateInstances(context.TODO(), terminateInput)
	if err != nil {
		log.LogError("Error happens when terminate instances %s : %s", strings.Join(instanceIDs, ","), err)
		return err
	} else {
		log.LogInfo("Terminate instances %s successfully", strings.Join(instanceIDs, ","))
	}
	if wait {
		err = client.WaitForInstanceTerminated(instanceIDs, timeout)
		if err != nil {
			log.LogError("Waiting for  instances %s termination timeout %s ", strings.Join(instanceIDs, ","), err)
			return err
		}

	}
	return nil
}

func (client *AwsV2Client) WaitForInstanceTerminated(instanceIDs []string, timeout time.Duration) error {
	log.LogInfo("Waiting for below instances terminated: %s ", strings.Join(instanceIDs, ","))
	_, err := client.WaitForInstancesTerminated(instanceIDs, timeout)
	return err
}

func (client *AwsV2Client) GetTagsOfInstanceProfile(instanceProfileName string) ([]iamtypes.Tag, error) {
	input := &iam.ListInstanceProfileTagsInput{
		InstanceProfileName: &instanceProfileName,
	}
	resp, err := client.iamClient.ListInstanceProfileTags(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	tags := resp.Tags
	return tags, err
}
