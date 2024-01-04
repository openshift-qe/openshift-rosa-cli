package aws_v1

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	awsclient "github.com/aws/aws-sdk-go/aws/client"

	awscredentials "github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/sts"
	CON "github.com/openshift-qe/openshift-rosa-cli/pkg/constants"
	"github.com/openshift-qe/openshift-rosa-cli/pkg/file"
)

// CreateAWSClient will return a valid AWS client for the userName
func CreateAWSClient(userName string, region string) (*AWSClient, *aws.Config, error) {
	// if IsValid(userName) {
	// 	return NewAWSClient(userName, region)
	// }
	if !CON.ENVCredential() && !CON.ENVAWSProfile() {
		if !IsValid(CON.DefaultAWSCredentialUser) {
			return nil, nil, fmt.Errorf("ERROR! Cannot create AWS session with %s credentials.Please Check the credential file%s",
				CON.DefaultAWSCredentialUser,
				CON.AWSCredentialsFilePath)
		}
	}
	defaultClient, cfg, err := NewAWSClient(CON.DefaultAWSCredentialUser, region)
	if userName == "" {
		return defaultClient, cfg, err
	}
	if err != nil {
		return nil, nil, fmt.Errorf("ERROR! Cannot create AWS session with %s credentials:%s",
			CON.DefaultAWSCredentialUser, err)
	}
	awsKeys, err := defaultClient.GrantNewAccessKeys(userName)
	if err != nil {
		return nil, nil, err
	}
	return NewAWSClient(userName, region, awsKeys)
}

func GrantValidAccessKeys(userName string) (*AccessKeyMod, error) {
	var keys awscredentials.Value
	var keysMod *AccessKeyMod
	var err error
	retryTimes := 3
	for retryTimes > 0 {
		if keys.AccessKeyID != "" {
			break
		}
		_, cfg, err := CreateAWSClient(userName, CON.DefaultAWSRegion)
		if err != nil {
			return nil, err
		}
		keys, err = cfg.Credentials.Get()
		fmt.Println(">>> Access key grant successfully")
		keysMod = &AccessKeyMod{
			AccessKeyId:     keys.AccessKeyID,
			SecretAccessKey: keys.SecretAccessKey,
		}
		retryTimes--
	}
	return keysMod, err
}

// NewAWSClient will return the AWS client with credentials, and the region is optional
func NewAWSClient(userName string, region string, awsKey ...*AccessKeyMod) (*AWSClient, *aws.Config, error) {
	var cfg *aws.Config
	AWSCredentialsFilePath := CON.AWSCredentialsFilePath
	if len(awsKey) == 1 {
		cfg = &aws.Config{
			Credentials: awscredentials.NewStaticCredentials(awsKey[0].AccessKeyId, awsKey[0].SecretAccessKey, ""),
			Retryer:     awsclient.DefaultRetryer{NumMaxRetries: 2},
			Region:      aws.String(region),
		}
	}
	if CON.ENVCredential() {
		cfg = &aws.Config{
			Credentials: awscredentials.NewEnvCredentials(),
			Retryer:     awsclient.DefaultRetryer{NumMaxRetries: 2},
			Region:      aws.String(region),
		}
	} else {
		if CON.ENVAWSProfile() {

			AWSCredentialsFilePath = os.Getenv("AWS_SHARED_CREDENTIALS_FILE")
		}
		cfg = &aws.Config{

			Credentials: awscredentials.NewSharedCredentials(AWSCredentialsFilePath, userName),
			Retryer:     awsclient.DefaultRetryer{NumMaxRetries: 2},
			Region:      aws.String(region),
		}

	}

	sess, err := session.NewSession(cfg)
	if err != nil {
		return nil, nil, err
	}
	return &AWSClient{
		ec2Client:            ec2.New(sess),
		route53Client:        route53.New(sess),
		stsClient:            sts.New(sess),
		iamClient:            iam.New(sess),
		cloudFormationClient: cloudformation.New(sess),
		kmsClient:            kms.New(sess),
		cloudWatchLogsClient: cloudwatchlogs.New(sess),
	}, cfg, nil
}

func IsValid(userName string) bool {
	if CON.ENVCredential() {
		client, _, err := NewAWSClient(userName, CON.DefaultAWSRegion)
		if err != nil {
			return false
		}
		valid, err := client.Valid()
		if err != nil || !valid {
			fmt.Println(err)
			return false
		}
		user, err := client.iamClient.GetUser(&iam.GetUserInput{})
		if err != nil || !valid {
			fmt.Println(err)
			return false
		}
		if *user.User.UserName == userName {
			return true
		}
		return false
	}
	iniCfg := file.IniConnection(CON.AWSCredentialsFilePath)
	var err error
	if _, exist := file.IsRecordExist(userName, iniCfg); !exist {
		fmt.Println(iniCfg.Sections())
		return false
	}
	client, _, err := NewAWSClient(userName, CON.DefaultAWSRegion)
	if err != nil {
		return false
	}
	valid, err := client.Valid()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return valid
}

func runCMD(cmd string) (stdout string, stderr string, err error) {
	var stdoutput bytes.Buffer
	var stderroutput bytes.Buffer
	CMD := exec.Command("bash", "-c", cmd)
	CMD.Stderr = &stderroutput
	CMD.Stdout = &stdoutput
	err = CMD.Run()
	stdout = strings.Trim(stdoutput.String(), "\n")
	stderr = strings.Trim(stderroutput.String(), "\n")
	return
}

func ShorterString(sourString string, length int) string {
	if len(sourString) <= length {
		return sourString
	}
	return sourString[:length]
}
