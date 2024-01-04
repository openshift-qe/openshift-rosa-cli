package aws_v1

import (
	"time"

	"github.com/aws/aws-sdk-go/service/iam"
)

func (client *AWSClient) UserExisted(userName string) bool {
	userInput := iam.GetUserInput{
		UserName: &userName,
	}
	_, err := client.iamClient.GetUser(&userInput)
	if err != nil {
		return false
	}
	return true
}

func (client *AWSClient) CreateUser(userName string, policyArn string) error {
	userInput := iam.CreateUserInput{
		PermissionsBoundary: &policyArn,
		UserName:            &userName,
	}
	_, err := client.iamClient.CreateUser(&userInput)
	return err
}

func (client *AWSClient) AttachPolicyToUser(userName string, policyArn string) error {
	policyAttach := iam.AttachUserPolicyInput{
		PolicyArn: &policyArn,
		UserName:  &userName,
	}
	_, err := client.iamClient.AttachUserPolicy(&policyAttach)
	timeout := 2
	start := 0

	for start < timeout {

		if attached, _ := client.PolicyAttachedToRole(userName, policyArn); attached {
			return nil
		}
		time.Sleep(1 * time.Minute)
		start++
	}
	return err
}

func (client *AWSClient) ListUserAttachedPolicies(userName string) ([]string, error) {
	policyArns := []string{}
	policyLister := iam.ListAttachedUserPoliciesInput{
		UserName: &userName,
	}
	policyOut, err := client.iamClient.ListAttachedUserPolicies(&policyLister)
	if err != nil {
		return policyArns, err
	}
	policies := policyOut.AttachedPolicies
	for _, policy := range policies {
		policyArns = append(policyArns, *policy.PolicyArn)
	}

	return policyArns, nil
}

func (client *AWSClient) PolicyAttachedToUser(userName string, policyArn string) (bool, error) {
	policies, err := client.ListUserAttachedPolicies(userName)
	if err != nil {
		return false, err
	}
	for _, policy := range policies {
		if policy == policyArn {
			return true, nil
		}
	}
	return false, nil
}
