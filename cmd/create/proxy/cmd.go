package proxy

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	vpcClient "gitlab.cee.redhat.com/openshift-group-I/ocm_aws/aws/vpc"
	awsV2 "gitlab.cee.redhat.com/openshift-group-I/ocm_aws/pkg/aws_client/aws_v2"
	. "gitlab.cee.redhat.com/openshift-group-I/ocm_aws/pkg/log"
)

var args struct {
	region      string
	vpcID       string
	zone        string
	imageID     string
	sshFilePath string
	caFilePath  string
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
		&args.vpcID,
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
		&args.sshFilePath,
		"ssh-file",
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
	Cmd.MarkFlagRequired("vpc-id")
	Cmd.MarkFlagRequired("region")
	Cmd.MarkFlagRequired("zone")
	Cmd.MarkFlagRequired("ca-file")
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
	_, ip, ca, err := vpc.LaunchProxyInstance(args.imageID, args.zone, args.sshFilePath)
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
