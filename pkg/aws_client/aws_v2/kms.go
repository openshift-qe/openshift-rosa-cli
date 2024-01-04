package aws_v2

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

func (client *AwsV2Client) CreateKMSKeys(tagKey string, tagValue string, description string, policy string, multiRegion bool) (keyID string, keyArn string, err error) {
	//Create the key

	result, err := client.kmsClient.CreateKey(context.TODO(), &kms.CreateKeyInput{
		Tags: []types.Tag{
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

func (client *AwsV2Client) DescribeKMSKeys(keyID string) (kms.DescribeKeyOutput, error) {
	// Create the key
	result, err := client.kmsClient.DescribeKey(context.TODO(), &kms.DescribeKeyInput{
		KeyId: &keyID,
	})
	if err != nil {
		fmt.Println("Got error describe key: ", err)
	}
	fmt.Println(*result)
	return *result, err
}
func (client *AwsV2Client) ScheduleKeyDeletion(kmsKeyId string, pendingWindowInDays int32) (*kms.ScheduleKeyDeletionOutput, error) {
	result, err := client.kmsClient.ScheduleKeyDeletion(context.TODO(), &kms.ScheduleKeyDeletionInput{
		KeyId:               aws.String(kmsKeyId),
		PendingWindowInDays: &pendingWindowInDays,
	})

	if err != nil {
		fmt.Println("Got error when ScheduleKeyDeletion: ", err)
	}

	// fmt.Println(*result.KeyMetadata.KeyId)
	return result, err
}

func (client *AwsV2Client) GetKMSPolicy(keyID string, policyName string) (kms.GetKeyPolicyOutput, error) {

	if policyName == "" {
		policyName = "default"
	}
	result, err := client.kmsClient.GetKeyPolicy(context.TODO(), &kms.GetKeyPolicyInput{
		KeyId:      &keyID,
		PolicyName: &policyName,
	})
	if err != nil {
		fmt.Println("Got error get KMS key policy: ", err)
	}
	return *result, err
}

func (client *AwsV2Client) PutKMSPolicy(keyID string, policyName string, policy string) (kms.PutKeyPolicyOutput, error) {
	if policyName == "" {
		policyName = "default"
	}
	result, err := client.kmsClient.PutKeyPolicy(context.TODO(), &kms.PutKeyPolicyInput{
		KeyId:      &keyID,
		PolicyName: &policyName,
		Policy:     &policy,
	})
	if err != nil {
		fmt.Println("Got error put KMS key policy: ", err)
	}
	return *result, err
}

func (client *AwsV2Client) TagKeys(kmsKeyId string, tagKey string, tagValue string) (*kms.TagResourceOutput, error) {

	output, err := client.kmsClient.TagResource(context.TODO(), &kms.TagResourceInput{
		KeyId: &kmsKeyId,
		Tags: []types.Tag{
			{
				TagKey:   aws.String(tagKey),
				TagValue: aws.String(tagValue),
			},
		},
	})
	if err != nil {
		fmt.Println("Got error add tag for KMS key: ", err)
	}
	return output, err
}
