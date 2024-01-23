package aws_v1

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/iam"
	CON "github.com/openshift-qe/openshift-rosa-cli/pkg/constants"
)

func (client *AWSClient) RoleExisted(roleName string) bool {
	_, err := client.GetRole(roleName)
	return err == nil
}

func (client *AWSClient) ListRoleAttachedPolicies(roleName string) ([]*iam.AttachedPolicy, error) {
	policies := []*iam.AttachedPolicy{}
	policyLister := iam.ListAttachedRolePoliciesInput{
		RoleName: &roleName,
	}
	policyOut, err := client.iamClient.ListAttachedRolePolicies(&policyLister)
	if err != nil {
		return policies, err
	}
	policies = policyOut.AttachedPolicies
	return policies, nil
}

func (client *AWSClient) DetachRolePolicies(roleName string) error {
	policies, err := client.ListRoleAttachedPolicies(roleName)
	if err != nil {
		return err
	}
	for _, policy := range policies {
		policyDetacher := iam.DetachRolePolicyInput{
			PolicyArn: policy.PolicyArn,
			RoleName:  &roleName,
		}
		_, err := client.iamClient.DetachRolePolicy(&policyDetacher)
		if err != nil {
			return err
		}
	}
	return nil
}

func (client *AWSClient) DeleteRoleInstanceProfiles(roleName string) error {
	inProfileLister := iam.ListInstanceProfilesForRoleInput{
		RoleName: &roleName,
	}
	out, err := client.iamClient.ListInstanceProfilesForRole(&inProfileLister)
	if err != nil {
		return err
	}
	for _, inProfile := range out.InstanceProfiles {
		profileDeleter := iam.RemoveRoleFromInstanceProfileInput{
			InstanceProfileName: inProfile.InstanceProfileName,
			RoleName:            &roleName,
		}
		_, err = client.iamClient.RemoveRoleFromInstanceProfile(&profileDeleter)
		if err != nil {
			return err
		}
	}

	return nil
}
func (client *AWSClient) ListInlinePolicys(roleName string) ([]*string, error) {
	listInput := iam.ListRolePoliciesInput{
		RoleName: &roleName,
	}
	output, err := client.iamClient.ListRolePolicies(&listInput)
	if err != nil {
		return nil, err
	}
	return output.PolicyNames, nil
}

func (client *AWSClient) DeleteInlinePolicys(roleName string) error {
	inLinePolicies, err := client.ListInlinePolicys(roleName)
	if err != nil {
		return err
	}
	for _, policy := range inLinePolicies {
		roleDeletion := iam.DeleteRolePolicyInput{
			PolicyName: policy,
			RoleName:   &roleName,
		}
		_, err = client.iamClient.DeleteRolePolicy(&roleDeletion)
		if err != nil {
			return err
		}
	}
	return nil
}

func (client *AWSClient) DeleteRole(roleName string) error {
	err := client.DetachRolePolicies(roleName)
	if err != nil {
		return err
	}
	err = client.DeleteInlinePolicys(roleName)
	if err != nil {
		return err
	}
	err = client.DeleteRoleInstanceProfiles(roleName)
	if err != nil {
		return err
	}
	roleDeleter := iam.DeleteRoleInput{
		RoleName: &roleName,
	}
	_, err = client.iamClient.DeleteRole(&roleDeleter)
	return err
}

func (client *AWSClient) CreateIAMRole(roleName string, policyArn string, externalID ...string) (string, error) {

	statement := map[string]interface{}{
		"Effect": "Allow",
		"Principal": map[string]interface{}{
			"Service": "ec2.amazonaws.com",
			"AWS": []string{
				CON.ProdENVTrustedRole,
				CON.StageENVTrustedRole,
				CON.StageIssuerTrustedRole,
			},
		},
		"Action": "sts:AssumeRole",
	}
	if len(externalID) == 1 {
		statement["Condition"] = map[string]map[string]string{
			"StringEquals": map[string]string{
				"sts:ExternalId": "aaaa",
			},
		}
	}
	return client.CreateRole(roleName, policyArn, statement)
}

func (client *AWSClient) CreateRegularRole(roleName string, policyArn string) (string, error) {

	statement := map[string]interface{}{
		"Effect": "Allow",
		"Principal": map[string]interface{}{
			"Service": "ec2.amazonaws.com",
		},
		"Action": "sts:AssumeRole",
	}
	return client.CreateRole(roleName, policyArn, statement)
}

func (client *AWSClient) CreateOperatorRole(roleName, policyArn string, idpArn string) (string, error) {
	conditionKeys := strings.Split(idpArn, "/")
	if len(conditionKeys) != 3 {
		err := fmt.Errorf(
			`Please check the idpArn: %s. \
			It should follow 
			arn:aws:iam::301721915996:oidc-provider/rh-oidc-staging.s3.us-east-1.amazonaws.com/1lfjp2tb7mqmuu6f2j08hvl3e8dc5mha`, idpArn)

		return "", err
	}
	conditionKey := strings.Join(conditionKeys[1:], "/")
	statement := map[string]interface{}{
		"Effect": "Allow",
		"Principal": map[string]interface{}{
			"Federated": idpArn,
		},
		"Action": "sts:AssumeRoleWithWebIdentity",
		"Condition": map[string]interface{}{
			"StringEquals": map[string]interface{}{
				fmt.Sprintf("%s:aud", conditionKey): "openshift",
			},
		},
	}
	return client.CreateRole(roleName, policyArn, statement)
}
func (client *AWSClient) CreateRole(roleName string, policyArn string, statements ...map[string]interface{}) (string, error) {
	if client.RoleExisted(roleName) {
		err := client.DeleteRole(roleName)
		if err != nil {
			err = fmt.Errorf("Cannot delete the existed role: %s due to error: %v", roleName, err)
			return "", err
		}
	}
	timeCreation := time.Now().Local().String()
	description := fmt.Sprintf("Created by OCM QE at %s", timeCreation)
	document := map[string]interface{}{
		"Version":   "2012-10-17",
		"Statement": []map[string]interface{}{},
	}
	if len(statements) != 0 {
		for _, statement := range statements {
			document["Statement"] = append(document["Statement"].([]map[string]interface{}), statement)
		}
	}
	documentBytes, err := json.Marshal(document)
	if err != nil {
		err = fmt.Errorf("Error to unmarshal the statement to string, %v", err)
	}
	documentStr := string(documentBytes)
	roleCreator := iam.CreateRoleInput{
		AssumeRolePolicyDocument: &documentStr,
		RoleName:                 &roleName,
		Description:              &description,
	}
	outRes, err := client.iamClient.CreateRole(&roleCreator)
	if err != nil {
		return "", err
	}
	roleArn := *outRes.Role.Arn
	err = client.AttachPolicy(roleName, policyArn)
	return roleArn, err
}

func (client *AWSClient) PolicyAttachedToRole(roleName string, policyArn string) (bool, error) {
	policies, err := client.ListRoleAttachedPolicies(roleName)
	if err != nil {
		return false, err
	}
	for _, policy := range policies {
		if *policy.PolicyArn == policyArn {
			return true, nil
		}
	}
	return false, nil
}

func (client *AWSClient) AttachPolicy(roleName string, policyArn string) error {
	policyAttach := iam.AttachRolePolicyInput{
		PolicyArn: &policyArn,
		RoleName:  &roleName,
	}
	_, err := client.iamClient.AttachRolePolicy(&policyAttach)
	if err != nil {
		return err
	}
	timeout := 2
	start := 0

	for start < timeout {

		if attached, _ := client.PolicyAttachedToRole(roleName, policyArn); attached {
			return nil
		}
		time.Sleep(1 * time.Minute)
		start++
	}
	return err
}

func (client *AWSClient) GetRole(roleName string) (role *iam.Role, err error) {
	roleGetter := iam.GetRoleInput{
		RoleName: &roleName,
	}
	roleout, err := client.iamClient.GetRole(&roleGetter)
	if err != nil {
		return
	}
	role = roleout.Role
	return

}

func (client *AWSClient) PolicyExisted(policyArn string) bool {
	_, err := client.GetPolicy(policyArn)
	return err == nil
}

func (client *AWSClient) GetPolicy(policyArn string) (policy *iam.Policy, err error) {
	policyGetter := iam.GetPolicyInput{
		PolicyArn: &policyArn,
	}
	policyout, err := client.iamClient.GetPolicy(&policyGetter)
	if err != nil {
		return
	}
	policy = policyout.Policy
	return

}
func (client *AWSClient) DeletePolicy(policyArn string) error {

	roleDeleter := iam.DeletePolicyInput{
		PolicyArn: &policyArn,
	}
	_, err := client.iamClient.DeletePolicy(&roleDeleter)
	return err
}

func (client *AWSClient) CreatePolicy(policyName string, statements ...map[string]interface{}) (string, error) {
	timeCreation := time.Now().Local().String()
	description := fmt.Sprintf("Created by OCM QE at %s", timeCreation)
	document := map[string]interface{}{
		"Version":   "2012-10-17",
		"Statement": []map[string]interface{}{},
	}
	if len(statements) != 0 {
		for _, statement := range statements {
			document["Statement"] = append(document["Statement"].([]map[string]interface{}), statement)
		}
	}
	documentBytes, err := json.Marshal(document)
	if err != nil {
		err = fmt.Errorf("Error to unmarshal the statement to string, %v", err)
	}
	documentStr := string(documentBytes)
	policyCreator := iam.CreatePolicyInput{
		PolicyDocument: &documentStr,
		PolicyName:     &policyName,
		Description:    &description,
	}
	outRes, err := client.iamClient.CreatePolicy(&policyCreator)
	if err != nil {
		return "", err
	}
	policyArn := *outRes.Policy.Arn
	return policyArn, err
}

func (client *AWSClient) CreateRoleForAuditLogForward(roleName, policyArn string, awsAccountID string, oidcEndpointURL string) (string, error) {
	statement := map[string]interface{}{
		"Effect": "Allow",
		"Principal": map[string]interface{}{
			"Federated": fmt.Sprintf("arn:aws:iam::%s:oidc-provider/%s", awsAccountID, oidcEndpointURL),
		},
		"Action": "sts:AssumeRoleWithWebIdentity",
		"Condition": map[string]interface{}{
			"StringEquals": map[string]interface{}{
				fmt.Sprintf("%s:sub", oidcEndpointURL): "system:serviceaccount:openshift-config-managed:cloudwatch-audit-exporter",
			},
		},
	}

	return client.CreateRole(roleName, policyArn, statement)
}

func (client *AWSClient) CreatePolicyForAuditLogForward(policyName string) (string, error) {

	statement := map[string]interface{}{
		"Effect":   "Allow",
		"Resource": "arn:aws:logs:*:*:*",
		"Action": []string{
			"logs:PutLogEvents",
			"logs:CreateLogGroup",
			"logs:PutRetentionPolicy",
			"logs:CreateLogStream",
			"logs:DescribeLogGroups",
			"logs:DescribeLogStreams",
		},
	}
	return client.CreatePolicy(policyName, statement)
}

func (client *AWSClient) DeleteRoleAndPolicy(roleName string) error {

	input := &iam.ListAttachedRolePoliciesInput{
		RoleName: &roleName,
	}
	output, err := client.iamClient.ListAttachedRolePolicies(input)
	if err != nil {
		return err
	}
	fmt.Println(output.AttachedPolicies)

	err = client.DeleteRole(roleName)
	if err != nil {
		return err
	}
	for _, policy := range output.AttachedPolicies {

		err = client.DeletePolicy(*policy.PolicyArn)
		if err != nil {
			return err
		}

	}
	return err
}
