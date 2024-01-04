package aws_v1

import (
	"crypto/sha1"
	"crypto/tls"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/service/iam"
)

func (client *AWSClient) IDPExisted(url string) (idpArn string, existed bool) {
	serverName := strings.TrimLeft(url, "https://")
	keyinput := iam.ListOpenIDConnectProvidersInput{}
	out, err := client.iamClient.ListOpenIDConnectProviders(&keyinput)
	if err != nil {
		return
	}
	providers := out.OpenIDConnectProviderList
	for _, IDP := range providers {
		if strings.Contains(*IDP.Arn, serverName) {
			idpArn = *IDP.Arn
			existed = true
		}
	}
	return

}

func (client *AWSClient) DeleteIDP(arn string) error {
	deleter := iam.DeleteOpenIDConnectProviderInput{
		OpenIDConnectProviderArn: &arn,
	}
	_, err := client.iamClient.DeleteOpenIDConnectProvider(&deleter)
	return err
}

// CreateOpenIDP will create an OpenIDP on AWS for sts cluster creation.
// The parameter url should be the oidc url of the cluster
func (client *AWSClient) CreateOpenIDP(url string) (idpArn string, err error) {
	// Check whether the OpenIDP Existing
	if arn, exsited := client.IDPExisted(url); exsited {
		err = client.DeleteIDP(arn)
		if err != nil {
			return
		}
	}

	// Create the new IDP
	OpenShift := "openshift"
	STSAWSCOM := "sts.amazonaws.com"

	thumbPrint, err := GetThumbPrint(url)
	if err != nil {
		return
	}

	idpCreator := iam.CreateOpenIDConnectProviderInput{
		ClientIDList: []*string{
			&OpenShift,
			&STSAWSCOM,
		},
		Url: &url,
		ThumbprintList: []*string{
			&thumbPrint,
		},
	}
	out, err := client.iamClient.CreateOpenIDConnectProvider(&idpCreator)
	if err != nil {
		return
	}
	idpArn = *out.OpenIDConnectProviderArn
	return
}

func GetThumbPrint(url string) (fingerPrint string, err error) {
	serverName := strings.TrimLeft(url, "https://")
	connectServer := fmt.Sprintf("%s:443", strings.Split(serverName, "/")[0])

	conf := tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", connectServer, &conf)
	if err != nil {
		return
	}
	certs := conn.ConnectionState().PeerCertificates
	fingerPrintBytes := sha1.Sum(certs[0].Raw)

	for _, finger := range fingerPrintBytes {
		hexString := strconv.FormatInt(int64(finger), 16)
		if len(hexString) == 1 {
			hexString = fmt.Sprintf("0%s", hexString)
		}
		fingerPrint = fingerPrint + strings.ToUpper(hexString)
	}
	return
}
