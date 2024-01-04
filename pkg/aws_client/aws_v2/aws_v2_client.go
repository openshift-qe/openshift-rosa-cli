package aws_v2

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	// elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	CON "gitlab.cee.redhat.com/openshift-group-I/ocm_aws/pkg/constants"
	"gitlab.cee.redhat.com/openshift-group-I/ocm_aws/pkg/log"
)

func CreateAWSV2Client(profileName string, region string) (*AwsV2Client, error) {
	var cfg aws.Config
	var err error

	if CON.ENVCredential() {
		log.LogInfo("Got AWS_ACCESS_KEY_ID env settings, going to build the config with the env")
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(region),
			config.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(
					os.Getenv("AWS_ACCESS_KEY_ID"),
					os.Getenv("AWS_SECRET_ACCESS_KEY"),
					"")),
		)
	} else {
		if CON.ENVAWSProfile() {
			file := os.Getenv("AWS_SHARED_CREDENTIALS_FILE")
			log.LogInfo("Got file path: %s from env variable AWS_SHARED_CREDENTIALS_FILE\n", file)
			cfg, err = config.LoadDefaultConfig(context.TODO(),
				config.WithRegion(region),
				config.WithSharedCredentialsFiles([]string{file}),
			)
		}
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(region),
			config.WithSharedConfigProfile(profileName),
		)
	}

	if err != nil {
		return nil, err
	}

	awsClient := &AwsV2Client{
		ec2Client:            ec2.NewFromConfig(cfg),
		route53Client:        route53.NewFromConfig(cfg),
		stackFormationClient: cloudformation.NewFromConfig(cfg),
		// elbClient:            elb.NewFromConfig(cfg),
		elbClient: elb.NewFromConfig(cfg),
		Region:    region,
		stsClient: sts.NewFromConfig(cfg),
		//stsClient:           sts.NewFromConfig(cfg),
		iamClient:     iam.NewFromConfig(cfg),
		clientContext: context.TODO(),
		kmsClient:     kms.NewFromConfig(cfg),

		//s3Client:            s3.NewFromConfig(cfg),
		//kmsClient:           kms.NewFromConfig(cfg),
		//servicequotasClient: servicequotas.NewFromConfig(cfg),
		//cloudWatchClient:    cloudwatch.NewFromConfig(cfg),
	}
	awsClient.AccountID = awsClient.GetAWSAccountID()
	return awsClient, nil
}
func (client *AwsV2Client) GetAWSAccountID() string {
	input := &sts.GetCallerIdentityInput{}
	out, err := client.stsClient.GetCallerIdentity(client.clientContext, input)
	if err != nil {
		return ""
	}
	return *out.Account
}

func (client *AwsV2Client) EC2() *ec2.Client {
	return client.ec2Client
}

func (client *AwsV2Client) Route53() *route53.Client {
	return client.route53Client
}
func (client *AwsV2Client) CloudFormation() *cloudformation.Client {
	return client.stackFormationClient
}
func (client *AwsV2Client) ELB() *elb.Client {
	return client.elbClient
}
