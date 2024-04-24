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

func (vpc *VPC) DeleteKeyPair(keyNames []string) error {
	for _, key := range keyNames {
		_, err := vpc.AWSClient.DeleteKeyPair(key)
		if err != nil {
			return err
		}
	}
	return nil
}
