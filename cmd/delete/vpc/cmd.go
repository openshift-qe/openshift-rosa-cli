package vpc

import (
	"github.com/spf13/cobra"
	vpcClient "gitlab.cee.redhat.com/openshift-group-I/ocm_aws/aws/vpc"
	awsV2 "gitlab.cee.redhat.com/openshift-group-I/ocm_aws/pkg/aws_client/aws_v2"
)

var args struct {
	region     string
	totalClean bool
	vpcID      string
}
var Cmd = &cobra.Command{
	Use:   "vpc",
	Short: "Delete vpc",
	Long:  "Delete vpc.",
	Example: `  # Delete a vpc with vpc ID
  ocmqe delete vpc --vpc-id <vpc id> --region us-east-2`,

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
		"id of the vpc",
	)
	flags.BoolVarP(
		&args.totalClean,
		"total-clean",
		"",
		false,
		"find the vpc with same name",
	)
	Cmd.MarkFlagRequired("vpc-id")
	Cmd.MarkFlagRequired("region")
}
func run(cmd *cobra.Command, _ []string) {
	console, err := awsV2.CreateAWSV2Client("", args.region)
	if err != nil {
		panic(err)
	}
	vpc, err := vpcClient.GenerateVPCByID(console, args.vpcID)
	if err != nil {
		panic(err)
	}
	vpc.DeleteVPCChain(args.totalClean)
}
