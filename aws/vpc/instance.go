package vpc

import (
	"strings"

	CON "gitlab.cee.redhat.com/openshift-group-I/ocm_aws/pkg/constants"
	"gitlab.cee.redhat.com/openshift-group-I/ocm_aws/pkg/log"
)

func (vpc *VPC) TerminateVPCInstances(nonClusterOnly bool) error {
	filters := []map[string][]string{
		{
			"vpc-id": []string{
				vpc.VpcID,
			},
		},
	}
	if nonClusterOnly {
		filters = append(filters, map[string][]string{
			"tag:Name": {
				CON.ProxyName,
				CON.BastionName,
			},
		})
	}
	insts, err := vpc.AWSClient.ListInstance(filters...)

	if err != nil {
		log.LogError("Error happened when list instances for vpc %s: %s", vpc.VpcID, err)
		return err
	}
	needTermination := []string{}
	for _, inst := range insts {
		needTermination = append(needTermination, *inst.InstanceId)
	}
	err = vpc.AWSClient.TerminateInstances(needTermination, true, 20)
	if err != nil {
		log.LogError("Terminating instances %s meet error: %s", strings.Join(needTermination, ","), err)
	} else {
		log.LogInfo("Terminating instances %s successfully", strings.Join(needTermination, ","))
	}
	return err

}