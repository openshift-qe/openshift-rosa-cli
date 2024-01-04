package aws_v1

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
)

func (client *AWSClient) CreateKMSKeys(tagKey string, tagValue string, description string, policy string, multiRegion bool) (string, string, error) {
	//Create the key
	result, err := client.kmsClient.CreateKey(&kms.CreateKeyInput{
		Tags: []*kms.Tag{
			{
				TagKey:   aws.String(tagKey),
				TagValue: aws.String(tagValue),
			},
		},
		Description: &description,
		Policy:      aws.String(policy),
		MultiRegion: &multiRegion,
	})

	if err != nil {
		fmt.Println("Got error creating key: ", err)
	}

	// fmt.Println(*result.KeyMetadata.KeyId)
	return *result.KeyMetadata.KeyId, *result.KeyMetadata.Arn, err
}

func (client *AWSClient) DescribeKMSKeys(keyID string) (kms.DescribeKeyOutput, error) {
	// Create the key
	result, err := client.kmsClient.DescribeKey(&kms.DescribeKeyInput{
		KeyId: &keyID,
	})
	if err != nil {
		fmt.Println("Got error describe key: ", err)
	}
	fmt.Println(*result)
	return *result, err
}
func (client *AWSClient) ScheduleKeyDeletion(kmsKeyId string, pendingWindowInDays int64) (*kms.ScheduleKeyDeletionOutput, error) {
	result, err := client.kmsClient.ScheduleKeyDeletion(&kms.ScheduleKeyDeletionInput{
		KeyId:               aws.String(kmsKeyId),
		PendingWindowInDays: &pendingWindowInDays,
	})

	if err != nil {
		fmt.Println("Got error when ScheduleKeyDeletion: ", err)
	}

	// fmt.Println(*result.KeyMetadata.KeyId)
	return result, err
}

func (client *AWSClient) GetKMSPolicy(keyID string, policyName string) (kms.GetKeyPolicyOutput, error) {

	if policyName == "" {
		policyName = "default"
	}
	result, err := client.kmsClient.GetKeyPolicy(&kms.GetKeyPolicyInput{
		KeyId:      &keyID,
		PolicyName: &policyName,
	})
	if err != nil {
		fmt.Println("Got error get KMS key policy: ", err)
	}
	return *result, err
}

func (client *AWSClient) PutKMSPolicy(keyID string, policyName string, policy string) (kms.PutKeyPolicyOutput, error) {
	if policyName == "" {
		policyName = "default"
	}
	result, err := client.kmsClient.PutKeyPolicy(&kms.PutKeyPolicyInput{
		KeyId:      &keyID,
		PolicyName: &policyName,
		Policy:     &policy,
	})
	if err != nil {
		fmt.Println("Got error put KMS key policy: ", err)
	}
	return *result, err
}
