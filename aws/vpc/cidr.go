package vpc

import (
	"fmt"
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
	CON "github.com/openshift-qe/openshift-rosa-cli/pkg/constants"
)

func NewCIDRPool(vpcCIDR string) *VPCCIDRPool {
	v := &VPCCIDRPool{
		CIDR: vpcCIDR,
	}
	prefix := CON.DefaultCIDRPrefix
	if v.Prefix != 0 {
		prefix = v.Prefix
	}
	v.GenerateSubnetPool(prefix)
	return v
}

func (v *VPCCIDRPool) GenerateSubnetPool(prefix int) {
	subnetcidrs := []*SubnetCIDR{}
	_, vpcSubnet, _ := net.ParseCIDR(v.CIDR)
	currentSubnet, finished := cidr.PreviousSubnet(vpcSubnet, prefix)
	for true {
		currentSubnet, finished = cidr.NextSubnet(currentSubnet, prefix)
		if !finished && vpcSubnet.Contains(currentSubnet.IP) {
			subnetcidr := SubnetCIDR{
				IPNet: currentSubnet,
				CIDR:  currentSubnet.String(),
			}
			subnetcidrs = append(subnetcidrs, &subnetcidr)
		} else {
			break
		}
	}
	v.SubNetPool = subnetcidrs
}

func (v *VPCCIDRPool) Allocate() *SubnetCIDR {
	for _, subnetCIDR := range v.SubNetPool {
		if !subnetCIDR.Reserved {
			subnetCIDR.Reserved = true
			return subnetCIDR
		}
	}
	return nil
}

// Reserve will reserve the ones you passed as parameter so you won't allocate them again from the pool
func (v *VPCCIDRPool) Reserve(reservedCIDRs ...string) error {
	for _, reservedCIDR := range reservedCIDRs {
		_, ipnet, err := net.ParseCIDR(reservedCIDR)
		if err != nil {
			return fmt.Errorf("you passed a wrong CIDR:%s for reserve. %s", reservedCIDR, err)
		}
		for _, freeCidr := range v.SubNetPool {
			if intersect(freeCidr.IPNet, ipnet) {
				freeCidr.Reserved = true
			}
		}
	}
	return nil
}

func intersect(n1, n2 *net.IPNet) bool {
	return n2.Contains(n1.IP) || n1.Contains(n2.IP)
}
