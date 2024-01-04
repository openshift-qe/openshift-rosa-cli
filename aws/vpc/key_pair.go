package vpc

import (
	"fmt"
)

func (vpc *VPC) CreateKeyPair(keyName string) (*string, error) {
	output, err := vpc.AWSClient.CreateKeyPair(keyName)
	if err != nil {
		return nil, err
	}
	fmt.Printf("create key pair: %v successfully\n", *output.KeyPairId)
	content := output.KeyMaterial

	return content, nil
}

func (vpc *VPC) DeleteKeyPair(keyName string) error {
	_, err := vpc.AWSClient.DeleteKeyPair(keyName)
	if err != nil {
		return err
	}
	fmt.Printf("delete key pair successfully\n")
	return nil
}
