package vpc

import (
	"fmt"
	"strings"

	awsV2 "gitlab.cee.redhat.com/openshift-group-I/ocm_aws/pkg/aws_client/aws_v2"
	CON "gitlab.cee.redhat.com/openshift-group-I/ocm_aws/pkg/constants"
	"gitlab.cee.redhat.com/openshift-group-I/ocm_aws/pkg/log"
)

//GenerateVPCByID will return a VPC with CIDRpool and subnets
// func (client *AwsV2Client) LoadFromAWS(vpcID string) (*VPC, error) {
// 	vpc, err := client.DescribeVPC(vpcID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	subnets, err := client.ListSubnetByVpcID(vpc.VpcID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	vpc.SubnetList = subnets
// 	reservedCIDRs := []string{}
// 	for _, sub := range subnets {
// 		reservedCIDRs = append(reservedCIDRs, sub.Cidr)
// 	}
// 	cidrPool := NewCIDRPool(vpc.CIDR)
// 	cidrPool.Reserve(reservedCIDRs...)
// 	vpc.CIDRPool = cidrPool
// 	return vpc, nil
// }

// GenerateVPCBySubnet will return a VPC with CIDRpool and subnets based on one of the subnet ID
func (vpc *VPC) GenerateVPCBySubnet(subnetID string) (*VPC, error) {
	log.LogInfo("Trying to load vpc from AWS by subnet: %s", subnetID)
	subnetDetail, err := vpc.AWSClient.ListSubnetDetail(subnetID)
	if err != nil {
		log.LogError("List subnet detail meets error: %s", err)
		return nil, err
	}
	log.LogInfo("Subnet info loaded from AWS by subnet: %s", subnetID)
	vpc, err = GenerateVPCByID(vpc.AWSClient, *subnetDetail[0].VpcId)
	log.LogInfo("VPC info loaded from AWS by subnet: %s", subnetID)
	return vpc, err
}

// CreateVPCChain create a complete set of web resources, including eip, vpc, subnet, route table, internet gateway, nat gateway, routes
// Inputs:
//
//	vpcCidr is a string of the vpc's cidr, e.g. "10.190.0.0/16".
//	region is a string of the AWS region. If this value is empty, the default region is "us-east-2".
//	zone is a slice. If only one subnet should be created, the first zone should be selected. If this value is empty, the default zone is "a".
//	If success, a VPC struct containing the ids of the created resources and nil.
//	Otherwise, nil and an error from the call.
func (vpc *VPC) CreateVPCChain(zones ...string) (*VPC, error) {
	log.LogInfo("Going to create vpc and the follow resources on zones: %s", strings.Join(zones, ","))
	respVpc, err := vpc.AWSClient.CreateVpc(vpc.CIDRValue, vpc.VPCName)
	if err != nil {
		log.LogError("Create vpc meets error: %s ", err.Error())
		return nil, err
	}
	log.LogInfo("VPC created on AWS with id: %s", *respVpc.Vpc.VpcId)
	_, err = vpc.AWSClient.ModifyVpcDnsAttribute(*respVpc.Vpc.VpcId, CON.VpcDnsHostnamesAttribute, true)
	if err != nil {
		log.LogError("Modify Vpc failed: %s ", err.Error())
		return nil, err
	}
	log.LogInfo("VPC DNS Updated on AWS with id: %s", *respVpc.Vpc.VpcId)
	vpc = vpc.ID(*respVpc.Vpc.VpcId)
	_, err = vpc.PrepareInternetGateway()
	if err != nil {
		log.LogError("Prepare Vpc internet gateway failed: %s ", err.Error())
		return vpc, err
	}
	log.LogInfo("Prepare vpc internetgateway for vpc %s", *respVpc.Vpc.VpcId)
	err = vpc.CreateMultiZoneSubnet(zones...)
	if err != nil {
		log.LogError("Create subnets meets error: %s", err.Error())
	} else {
		log.LogInfo("Create subnets successfully")
	}
	return vpc, err
}

func (vpc *VPC) DeleteVPCChain(totalClean ...bool) error {
	vpcID := vpc.VpcID
	if vpcID == "" {
		return fmt.Errorf("got empty vpc ID to clean. Make sure you loaded it from AWS")
	}
	log.LogInfo("Going to delete the vpc and follow resources by ID: %s", vpcID)
	log.LogInfo("Going to terminate proxy instances if existing")
	err := vpc.TerminateVPCInstances(true)
	if err != nil {
		log.LogError("Delete vpc instances meets error: %s", err.Error())
		return err
	}
	log.LogInfo("Delete vpc instances successfully")

	log.LogInfo("Going to delete proxy security group")
	err = vpc.DeleteVPCSecurityGroups(true)
	if err != nil {
		log.LogError("Delete vpc proxy security group meets error: %s", err.Error())
		return err
	}
	log.LogInfo("Delete vpc vpc proxy security group successfully")

	err = vpc.DeleteVPCRouteTables(vpcID)
	if err != nil {
		log.LogError("Delete vpc route tables meets error: %s", err.Error())
		return err
	}
	log.LogInfo("Delete vpc route tables successfully")

	err = vpc.DeleteVPCNatGateways(vpcID)
	if err != nil {
		log.LogError("Delete vpc nat gatways meets error: %s", err.Error())
		return err
	}

	log.LogInfo("Delete vpc nat gateways successfully")
	err = vpc.DeleteVPCEndpoints()
	if err != nil {
		log.LogError("Delete vpc endpoints meets error: %s", err.Error())
		return err
	}
	if len(totalClean) == 1 && totalClean[0] {
		log.LogInfo("Got total clean set, going to delete other possible resource leak")
		err = vpc.DeleteVPCELBs()
		if err != nil {
			log.LogError("Delete vpc load balancers meets error: %s", err.Error())
			return err
		}
		err = vpc.DeleteVPCSecurityGroups(false)
		if err != nil {
			log.LogError("Delete vpc security groups meets error: %s", err.Error())
			return err
		}
	}
	err = vpc.DeleteVPCNetworkInterfaces()
	if err != nil {
		log.LogError("Delete vpc network interfaces meets error: %s", err.Error())
		return err
	}

	err = vpc.DeleteVPCInternetGateWays()
	if err != nil {
		log.LogError("Delete vpc internet gatways meets error: %s", err.Error())
		return err
	}

	err = vpc.DeleteVPCSubnets()
	if err != nil {
		log.LogError("Delete vpc subnets meets error: %s", err.Error())
		return err
	}

	_, err = vpc.AWSClient.DeleteVpc(vpc.VpcID)
	if err != nil {
		log.LogError("Delete vpc  meets error: %s", err.Error())
		return err
	}
	return nil
}

// PrepareVPC will find a vpc named <vpcName>
// If there is no vpc in the name
// It will Create vpc with the name in the region
// checkExisting means if you want to check current existing vpc to re-use.
// Just be careful once you use checkExisting, the vpc may have subnets not existing in your zones. And maybe multi subnets in the zones
// Try vpc.PreparePairSubnets by zone for further implementation to get a pair of
// Zones will be customized if you want. Otherwise, it will use the default zone "a"
func PrepareVPC(vpcName string, region string, vpcCIDR string, checkExisting bool, zones ...string) (*VPC, error) {
	if vpcCIDR == "" {
		vpcCIDR = CON.DefaultVPCCIDR
	}
	log.LogInfo("Going to prepare a vpc with name %s, on region %s, with cidr %s and subnets on zones %s",
		vpcName, region, vpcCIDR, strings.Join(zones, ","))
	awsclient, err := awsV2.CreateAWSV2Client("", region)
	if err != nil {
		log.LogError("Create AWS Client due to error: %s", err.Error())
		return nil, err
	}
	if checkExisting {
		log.LogInfo("Got checkExisting set to true, will check if there is existing vpc in same name")
		vpcs, err := awsclient.ListVPCByName(vpcName)
		if err != nil {
			log.LogError("Error happened when try to find a vpc: %s", err.Error())
			return nil, err
		}
		if len(vpcs) != 0 {
			vpcID := *vpcs[0].VpcId
			log.LogInfo("Got a vpc %s with name %s on region %s. Just load it for usage",
				vpcID, vpcName, region)
			vpc, err := GenerateVPCByID(awsclient, vpcID)
			if err != nil {
				log.LogError("Load vpc %s details meets error %s",
					vpcID, err.Error())
				return nil, err
			}
			for _, zone := range zones {
				_, err = vpc.PreparePairSubnetByZone(zone)
				if err != nil {
					log.LogError("Prepare subnets for vpc %s on zone %s meets error %s",
						vpcID, zone, err.Error())
					return nil, err
				}
			}
			return vpc, nil
		}
		log.LogInfo("Got no vpc with name %s on region %s. Going to create a new one",
			vpcName, region)
	}

	vpc := NewVPC().
		Name(vpcName).
		AWSclient(awsclient).
		SetRegion(region).
		CIDR(vpcCIDR).
		NewCIDRPool()
	vpc, err = vpc.CreateVPCChain(zones...)
	if err != nil {
		log.LogError("Create vpc chain meets error: %s", err.Error())
	} else {
		log.LogInfo("Create vpc chain successfully. Enjoy it.")
	}

	return vpc, err
}

// NewVPC will return a new VPC instance
// CIDR can be empty, then it will use default value

// GenerateVPCByID will return a VPC with CIDRpool and subnets
// If you know the vpc ID on AWS, then try to generate it
func GenerateVPCByID(awsClient *awsV2.AwsV2Client, vpcID string) (*VPC, error) {
	vpc := NewVPC().AWSclient(awsClient).ID(vpcID)
	vpcResp, err := vpc.AWSClient.DescribeVPC(vpcID)
	if err != nil {
		return nil, err
	}
	vpc = vpc.Name(getTagName((vpcResp.Tags))).SetRegion(awsClient.Region).CIDR(*vpcResp.CidrBlock)
	if err != nil {
		return nil, err
	}
	_, err = vpc.ListSubnets()
	if err != nil {
		return nil, err
	}
	reservedCIDRs := []string{}
	for _, sub := range vpc.SubnetList {
		reservedCIDRs = append(reservedCIDRs, sub.Cidr)
	}
	cidrPool := NewCIDRPool(vpc.CIDRValue)
	cidrPool.Reserve(reservedCIDRs...)
	vpc.CIDRPool = cidrPool
	return vpc, nil
}

// GenerateVPCBySubnet will return a VPC with CIDRpool and subnets based on one of the subnet ID
// If you know the subnet ID on AWS, then try to generate it on AWS.
func GenerateVPCBySubnet(awsClient *awsV2.AwsV2Client, subnetID string) (*VPC, error) {
	subnetDetail, err := awsClient.ListSubnetDetail(subnetID)
	if err != nil {
		return nil, err
	}
	vpc, err := GenerateVPCByID(awsClient, *subnetDetail[0].VpcId)
	return vpc, err
}
