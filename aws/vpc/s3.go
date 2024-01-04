package vpc

func (vpc *VPC) DeleteVPCEndpoints() error {
	err := vpc.AWSClient.DeleteVPCEndpoints(vpc.VpcID)

	return err
}
