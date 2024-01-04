package aws_v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	CON "github.com/openshift-qe/openshift-rosa-cli/pkg/constants"
	"github.com/openshift-qe/openshift-rosa-cli/pkg/file"
)

// IsValid will return whether the credential is valid
// If userName not existed in .aws/credentials file it will return false
// if the credential in the .aws/credentials file is not valid, it will return false

func (client *AWSClient) Valid() (bool, error) {
	_, err := client.stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		switch typedErr := err.(type) {
		case awserr.RequestFailure:
			return client.handleAWSRequestFailure(typedErr)
		default:
			return false, err
		}
	}
	return true, nil
}

func (client *AWSClient) handleAWSRequestFailure(err awserr.RequestFailure) (bool, error) {
	switch err.StatusCode() {
	case http.StatusForbidden:
		return false, nil
	// An error occurred trying to check validity of credentials.
	case http.StatusBadRequest:
		return false, nil
	default:
		return false, err
	}
}

// GrantNewAccessKeys will return new access key
// If access keys reached the limited max number, the latest one will be deleted
func (client *AWSClient) GrantNewAccessKeys(IAMUserName string) (*AccessKeyMod, error) {
	fmt.Println(">>> Going to get new credentials")
	var outPut *iam.CreateAccessKeyOutput
	var err error
	if !client.UserExisted(IAMUserName) {
		err := client.CreateUser(IAMUserName, CON.AdminPolicyArn)
		if err != nil {
			return nil, fmt.Errorf("User %s creation failed with error: %v", IAMUserName, err)
		}
		err = client.AttachPolicyToUser(IAMUserName, CON.AdminPolicyArn)
		if err != nil {
			return nil, fmt.Errorf("Attach policy %s to User %s failed with error: %v", CON.AdminPolicyArn, IAMUserName, err)
		}
	}
	outPut, err = client.iamClient.CreateAccessKey(&iam.CreateAccessKeyInput{UserName: &IAMUserName})
	if err != nil {
		switch typedErr := err.(type) {
		case awserr.RequestFailure:
			if typedErr.StatusCode() == CON.HTTPConflict {

				laterAccessKey, err := client.GetLaterAccessKeyID(IAMUserName)
				if err != nil {
					return nil, err
				}
				err = client.DeleteAccessKey(IAMUserName, laterAccessKey)
				if err != nil {
					return nil, err
				}
				outPut, err = client.iamClient.CreateAccessKey(&iam.CreateAccessKeyInput{UserName: &IAMUserName})
				if err != nil {
					return nil, err
				}
			}
		default:
			return nil, err
		}
	}
	time.Sleep(time.Duration(30 * time.Second))
	accessKey := outPut.AccessKey
	if accessKey == nil {
		return nil, fmt.Errorf("ERROR! Got empty access keys")
	}
	awsKeymod, err := RecordAWSAccessKey(IAMUserName, accessKey)
	return awsKeymod, err
}

// DeleteAccessKey will delete the <accessKey> owned by <userName>
func (client *AWSClient) DeleteAccessKey(userName string, accessKeyID string) error {
	_, err := client.iamClient.DeleteAccessKey(&iam.DeleteAccessKeyInput{AccessKeyId: &accessKeyID, UserName: &userName})
	return err
}

// GetLaterAccessKeyID will return the later AccessKey ID in the list
func (client *AWSClient) GetLaterAccessKeyID(IAMUserName string) (string, error) {
	listOutPut, err := client.iamClient.ListAccessKeys(&iam.ListAccessKeysInput{UserName: &IAMUserName})
	if err != nil {
		return "", err
	}
	accessKeyMetadatas := listOutPut.AccessKeyMetadata
	key := accessKeyMetadatas[0]
	for _, metadata := range accessKeyMetadatas {
		if metadata.CreateDate.After(*key.CreateDate) {
			key = metadata
		}

	}
	return *key.AccessKeyId, nil
}

// RecordAWSAccessKey will record the access keys of the userName to the file defined in constants
func RecordAWSAccessKey(userName string, accessKey *iam.AccessKey) (*AccessKeyMod, error) {
	fmt.Println(">>> Going to record access keys")
	accessKeyMod := &AccessKeyMod{
		AccessKeyId:     *accessKey.AccessKeyId,
		SecretAccessKey: *accessKey.SecretAccessKey,
	}
	if CON.ENVCredential() {
		return accessKeyMod, nil
	}
	iniCfg := file.IniConnection(CON.AWSCredentialsFilePath)
	iniCfg.BlockMode = true
	var err error
	if _, exist := file.IsRecordExist(userName, iniCfg); exist {
		sec, _ := iniCfg.GetSection(userName)
		err = sec.ReflectFrom(accessKeyMod)
	} else {
		sec, _ := iniCfg.NewSection(userName)
		err = sec.ReflectFrom(accessKeyMod)
	}

	if err == nil {
		err = iniCfg.SaveTo(CON.AWSCredentialsFilePath)
	}

	if err != nil {
		fmt.Printf(">>>[INI][Warning] Writing AWS Access Key into file meet error: %v\n", err)
	}
	return accessKeyMod, err
}
