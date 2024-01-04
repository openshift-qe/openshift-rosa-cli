package vpc

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	CON "github.com/openshift-qe/openshift-rosa-cli/pkg/constants"
	"github.com/openshift-qe/openshift-rosa-cli/pkg/log"
)

func (vpc *VPC) DeleteVPCSecurityGroups(customizedOnly bool) error {
	needCleanGroups := []types.SecurityGroup{}
	securityGroups, err := vpc.AWSClient.ListSecurityGroups(vpc.VpcID)
	if customizedOnly {
		for _, sg := range securityGroups {
			for _, tag := range sg.Tags {
				if *tag.Key == "Name" && (*tag.Value == CON.ProxySecurityGroupName || *tag.Value == CON.AdditionalSecurityGroupName) {
					needCleanGroups = append(needCleanGroups, sg)
				}
			}
		}
	} else {
		needCleanGroups = securityGroups
	}
	if err != nil {
		return err
	}
	for _, sg := range needCleanGroups {
		_, err = vpc.AWSClient.DeleteSecurityGroup(*sg.GroupId)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateAndAuthorizeDefaultSecurityGroupForMS Create security group and Authorize ingress for managed service testing.
// Ports aws.DefaultSGIngressPorts is defined for create 'odf-sec-group'
func (vpc *VPC) CreateAndAuthorizeDefaultSecurityGroupForMS(cidr string, protocol string, ports []map[string]int32) (string, error) {
	if protocol == "" {
		protocol = CON.TCPProtocol
	}
	resp, err := vpc.AWSClient.CreateSecurityGroup(vpc.VpcID, CON.SecurityGroupName, CON.SecurityGroupDescription)
	if err != nil {
		log.LogError("Create security group failed" + *resp.GroupId)
		return "", err
	}
	for _, i := range ports {
		fmt.Println(i["fromPort"])
		vpc.AWSClient.AuthorizeSecurityGroupIngress(*resp.GroupId, cidr, protocol, i["fromPort"], i["toPort"])
	}
	return *resp.GroupId, err
}

func (vpc *VPC) CreateAndAuthorizeDefaultSecurityGroupForProxy() (string, error) {
	var groupID string
	var err error
	protocol := CON.TCPProtocol
	resp, err := vpc.AWSClient.CreateSecurityGroup(vpc.VpcID, CON.ProxySecurityGroupName, CON.ProxySecurityGroupDescription)
	if err != nil {
		log.LogError("Create proxy security group failed for vpc %s: %s", vpc.VpcID, err)
		return "", err
	}
	groupID = *resp.GroupId
	log.LogInfo("SG %s created for vpc %s", groupID, vpc.VpcID)
	cidrPortsMap := map[string]int32{
		vpc.CIDRValue: 8080,
		"0.0.0.0/0":   22,
	}
	for cidr, port := range cidrPortsMap {
		_, err = vpc.AWSClient.AuthorizeSecurityGroupIngress(groupID, cidr, protocol, port, port)
		if err != nil {
			log.LogError("Authorize CIDR %s with port %s failed to SG %s of vpc %s: %s",
				cidr, port, groupID, vpc.VpcID, err)
			return groupID, err
		}
	}
	log.LogInfo("Authorize SG %s successfully for proxy.", groupID)

	return groupID, err
}
