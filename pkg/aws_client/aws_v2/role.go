package aws_v2

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

func (client *AwsV2Client) CreateRole(roleName string,
	assumeRolePolicyDocument string,
	permissionBoundry string,
	tags map[string]string,
	path string,
) (types.Role, error) {
	var roleTags []types.Tag
	for tagKey, tagValue := range tags {
		roleTags = append(roleTags, types.Tag{
			Key:   &tagKey,
			Value: &tagValue,
		})
	}
	description := "This is created role for ocm-qe automation testing"
	input := &iam.CreateRoleInput{
		RoleName:                 &roleName,
		AssumeRolePolicyDocument: &assumeRolePolicyDocument,
		Path:                     &path,
		PermissionsBoundary:      &permissionBoundry,
		Tags:                     roleTags,
		Description:              &description,
	}
	resp, err := client.iamClient.CreateRole(context.TODO(), input)
	if err != nil {
		return *resp.Role, err
	}
	err = client.WaitForResourceExisting("role-"+*resp.Role.RoleName, 10) // add a prefix to meet the resourceExisting split rule
	return *resp.Role, err
}

func (client *AwsV2Client) GetRole(roleName string) (*types.Role, error) {
	input := &iam.GetRoleInput{
		RoleName: &roleName,
	}
	out, err := client.iamClient.GetRole(context.TODO(), input)
	return out.Role, err
}
func (client *AwsV2Client) DeleteRole(roleName string) error {

	input := &iam.DeleteRoleInput{
		RoleName: &roleName,
	}
	_, err := client.iamClient.DeleteRole(context.TODO(), input)
	return err
}

func (client *AwsV2Client) DeleteRoleAndPolicy(roleName string, managedPolicy bool) error {
	input := &iam.ListAttachedRolePoliciesInput{
		RoleName: &roleName,
	}
	output, err := client.iamClient.ListAttachedRolePolicies(client.clientContext, input)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	fmt.Println(output.AttachedPolicies)
	for _, policy := range output.AttachedPolicies {
		err = client.DetachIAMPolicy(roleName, *policy.PolicyArn)
		if err != nil {
			return err
		}
		if !managedPolicy {
			err = client.DeletePolicy(*policy.PolicyArn)
			if err != nil {
				return err
			}
		}

	}
	err = client.DeleteRole(roleName)
	return err
}

func (client *AwsV2Client) ListRoles() ([]types.Role, error) {
	input := &iam.ListRolesInput{}
	out, err := client.iamClient.ListRoles(context.TODO(), input)
	return out.Roles, err
}

// func (client *AwsV2Client)FilterNeedCleanRoles(duringDate int)([]types.Role, error) {

// 	allroles , err := client.ListRoles()
// 	if err !=nil{
// 		return nil, err
// 	}
// 	for _, role := range allroles{
// 		if role.CreateDate
// 	}
// }
