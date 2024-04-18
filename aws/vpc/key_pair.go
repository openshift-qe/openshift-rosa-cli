package vpc

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func (vpc *VPC) CreateKeyPair(keyName string) (*ec2.CreateKeyPairOutput, error) {
	output, err := vpc.AWSClient.CreateKeyPair(keyName)
	if err != nil {
		return nil, err
	}
	fmt.Printf("create key pair: %v successfully\n", *output.KeyPairId)

	return output, nil
}

func (vpc *VPC) DeleteKeyPair(keyName string) error {
	_, err := vpc.AWSClient.DeleteKeyPair(keyName)
	if err != nil {
		return err
	}
	fmt.Printf("delete key pair successfully\n")
	return nil
}
