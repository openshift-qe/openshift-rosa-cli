package aws_v1

import (
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

type AccessKeyMod struct {
	AccessKeyId     string `ini:"aws_access_key_id,omitempty"`
	SecretAccessKey string `ini:"aws_secret_access_key,omitempty"`
}
type AWSClient struct {
	ec2Client            ec2iface.EC2API
	route53Client        route53iface.Route53API
	stsClient            stsiface.STSAPI
	iamClient            iamiface.IAMAPI
	cloudFormationClient cloudformationiface.CloudFormationAPI
	kmsClient            kmsiface.KMSAPI
	cloudWatchLogsClient cloudwatchlogsiface.CloudWatchLogsAPI
}
