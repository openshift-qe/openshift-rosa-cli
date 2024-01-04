package aws_v2

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

func (client *AwsV2Client) CreateIAMPolicy(policyName string, policyDocument string, tags map[string]string) (*types.Policy, error) {
	var policyTags []types.Tag
	for tagKey, tagValue := range tags {
		policyTags = append(policyTags, types.Tag{
			Key:   &tagKey,
			Value: &tagValue,
		})
	}
	description := "Policy for ocm-qe testing"
	input := &iam.CreatePolicyInput{
		PolicyName:     &policyName,
		PolicyDocument: &policyDocument,
		Tags:           policyTags,
		Description:    &description,
	}
	output, err := client.iamClient.CreatePolicy(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	err = client.WaitForResourceExisting("policy-"+*output.Policy.Arn, 10) // add a prefix to meet the resourceExisting split rule
	return output.Policy, err
}

func (client *AwsV2Client) GetIAMPolicy(policyArn string) (*types.Policy, error) {
	input := &iam.GetPolicyInput{
		PolicyArn: &policyArn,
	}
	out, err := client.iamClient.GetPolicy(context.TODO(), input)
	// out.Policy.Tags
	return out.Policy, err
}

func (client *AwsV2Client) DeleteIAMPolicy(arn string) error {
	input := &iam.DeletePolicyInput{
		PolicyArn: &arn,
	}
	err := client.DeletePolicyVersions(arn)
	if err != nil {
		return err
	}
	_, err = client.iamClient.DeletePolicy(context.TODO(), input)
	return err
}
func (client *AwsV2Client) ListIAMPolicy(prefix string) {
	// input := iam.ListPoliciesInput
}
func (client *AwsV2Client) AttachIAMPolicy(roleName string, policyArn string) error {
	input := &iam.AttachRolePolicyInput{
		PolicyArn: &policyArn,
		RoleName:  &roleName,
	}
	_, err := client.iamClient.AttachRolePolicy(context.TODO(), input)
	return err

}
func (client *AwsV2Client) DetachIAMPolicy(roleAName string, policyArn string) error {
	input := &iam.DetachRolePolicyInput{
		RoleName:  &roleAName,
		PolicyArn: &policyArn,
	}
	_, err := client.iamClient.DetachRolePolicy(context.TODO(), input)
	return err
}
func (client *AwsV2Client) GetCustomerIAMPolicies() ([]types.Policy, error) {

	maxItem := int32(1000)
	input := &iam.ListPoliciesInput{
		Scope:    "Local",
		MaxItems: &maxItem,
	}
	out, err := client.iamClient.ListPolicies(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	// for _, policy := range out.Policies {
	// 	fmt.Println(*policy.Arn)
	// 	fmt.Println(policy.CreateDate.Date())
	// 	fmt.Println(*policy.AttachmentCount)
	// }
	return out.Policies, err

}
func CleanByOutDate(policy types.Policy) bool {
	now := time.Now().UTC()
	return policy.CreateDate.Add(7 * time.Hour * 24).Before(now)
}

func CleanByName(policy types.Policy) bool {
	return strings.Contains(*policy.PolicyName, "sdq-ci-")
}

func (client *AwsV2Client) FilterNeedCleanPolicies(cleanRule func(types.Policy) bool) ([]types.Policy, error) {
	needClean := []types.Policy{}

	policies, err := client.GetCustomerIAMPolicies()
	if err != nil {
		return needClean, err
	}
	for _, policy := range policies {
		// fmt.Printf("policy %s with creation date %s with attachment %d\n", *policy.Arn, policy.CreateDate.GoString(), *policy.AttachmentCount)
		if cleanRule(policy) {

			needClean = append(needClean, policy)
		}
	}
	return needClean, nil
}

func (client *AwsV2Client) DeletePolicy(arn string) error {
	input := &iam.DeletePolicyInput{
		PolicyArn: &arn,
	}
	err := client.DeletePolicyVersions(arn)
	if err != nil {
		return err
	}
	_, err = client.iamClient.DeletePolicy(context.TODO(), input)
	return err
}

func (client *AwsV2Client) DeletePolicyVersions(policyArn string) error {
	input := &iam.ListPolicyVersionsInput{
		PolicyArn: &policyArn,
	}
	out, err := client.iamClient.ListPolicyVersions(context.TODO(), input)
	if err != nil {
		return err
	}
	for _, version := range out.Versions {
		if version.IsDefaultVersion {
			continue
		}
		input := &iam.DeletePolicyVersionInput{
			PolicyArn: &policyArn,
			VersionId: version.VersionId,
		}
		_, err = client.iamClient.DeletePolicyVersion(context.TODO(), input)
		if err != nil {
			return err
		}
	}
	return nil
}
func (client *AwsV2Client) CleanPolicies(cleanRule func(types.Policy) bool) error {
	policies, err := client.FilterNeedCleanPolicies(cleanRule)
	if err != nil {
		return err
	}
	for _, policy := range policies {
		if *policy.AttachmentCount == 0 {
			fmt.Println("Can be deleted: ", *policy.Arn)
			err = client.DeletePolicy(*policy.Arn)

			if err != nil {
				return err
			}
		} else {
			// fmt.Printf("Policy %s has role attachment\n", *policy.Arn)
		}
	}
	return nil
}
