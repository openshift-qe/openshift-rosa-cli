package vpc

import (
	"sync"
)

func (vpc *VPC) DeleteVPCNatGateways(vpcID string) error {
	var delERR error
	natGateways, err := vpc.AWSClient.ListNatGateWays(vpcID)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, ngw := range natGateways {
		wg.Add(1)
		go func(gateWayID string) {
			defer wg.Done()
			_, err = vpc.AWSClient.DeleteNatGateway(gateWayID, 120)
			if err != nil {
				delERR = err
			}
		}(*ngw.NatGatewayId)
	}
	wg.Wait()
	return delERR
}
