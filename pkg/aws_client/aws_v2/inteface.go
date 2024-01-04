package aws_v2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	// elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/route53"
)

type AwsV2Client struct {
	ec2Client            *ec2.Client
	route53Client        *route53.Client
	stackFormationClient *cloudformation.Client
	elbClient            *elb.Client
	stsClient            *sts.Client
	Region               string
	//stsClient           *sts.Client
	iamClient     *iam.Client
	clientContext context.Context
	AccountID     string
	//s3Client            *s3.Client
	kmsClient *kms.Client
	//servicequotasClient *servicequotas.Client
	//cloudWatchClient    *cloudwatch.Client
}
