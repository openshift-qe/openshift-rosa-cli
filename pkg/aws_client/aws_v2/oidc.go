package aws_v2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iam"
)

func (client *AwsV2Client) DeleteOIDCProvider(providerArn string) error {
	input := &iam.DeleteOpenIDConnectProviderInput{
		OpenIDConnectProviderArn: &providerArn,
	}
	_, err := client.iamClient.DeleteOpenIDConnectProvider(context.TODO(), input)
	return err
}
