package proxy

import (
	"fmt"
	"os"

	vpcClient "github.com/openshift-qe/openshift-rosa-cli/aws/vpc"
	awsV2 "github.com/openshift-qe/openshift-rosa-cli/pkg/aws_client/aws_v2"
	. "github.com/openshift-qe/openshift-rosa-cli/pkg/log"
	"github.com/spf13/cobra"
)

var args struct {
	region         string
	vpcID          string
	zone           string
	imageID        string
	privateKeyPath string
	keyPairName    string
	caFilePath     string
}

var Cmd = &cobra.Command{
	Use:   "proxy",
	Short: "Create proxy",
	Long:  "Create proxy.",
	Example: `  # Create a proxy
  ocmqe create proxy  --region us-east-2 --vpc-id <vpc id>`,

	Run: run,
}

func init() {
	flags := Cmd.Flags()
	flags.SortFlags = false
	flags.StringVarP(
		&args.region,
		"region",
		"",
		"",
		"Region of the vpc",
	)
	flags.StringVarP(
		&args.vpcID,
		"vpc-id",
		"",
		"",
		"Create a pair of subnets to",
	)
	flags.StringVarP(
		&args.zone,
		"zone",
		"",
		"",
		"Create a proxy to the indicated zone",
	)
	flags.StringVarP(
		&args.imageID,
		"image-id",
		"",
		"",
		"Create a proxy to the indicated zone",
	)

	flags.StringVarP(
		&args.caFilePath,
		"ca-file",
		"",
		"",
		"Create a proxy and store the ca file",
	)

	flags.StringVarP(
		&args.keyPairName,
		"keypair-name",
		"",
		"",
		"Store key pair private key to the path",
	)
	flags.StringVarP(
		&args.privateKeyPath,
		"privatekey-path",
		"",
		"",
		"Store key pair private key to the path",
	)
	Cmd.MarkFlagRequired("vpc-id")
	Cmd.MarkFlagRequired("region")
	Cmd.MarkFlagRequired("zone")
	Cmd.MarkFlagRequired("ca-file")
	Cmd.MarkFlagRequired("keypair-name")
	Cmd.MarkFlagRequired("privatekey-path")
}
func run(cmd *cobra.Command, _ []string) {
	console, err := awsV2.CreateAWSV2Client("", "us-west-2")
	if err != nil {
		panic(err)
	}
	vpc, err := vpcClient.GenerateVPCByID(console, args.vpcID)
	if err != nil {
		panic(err)
	}
	_, ip, ca, err := vpc.LaunchProxyInstance(args.imageID, args.zone, args.keyPairName, args.privateKeyPath)
	if err != nil {
		panic(err)
	}
	httpProxy := fmt.Sprintf("http://%s:8080", ip)
	httpsProxy := fmt.Sprintf("https://%s:8080", ip)
	file, err := os.OpenFile(args.caFilePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	_, err = file.WriteString(ca)
	if err != nil {
		panic(err)
	}
	LogInfo("HTTP PROXY: %s", httpProxy)
	LogInfo("HTTPs PROXY: %s", httpsProxy)
	LogInfo("CA FILE PATH: %s", args.caFilePath)
}
