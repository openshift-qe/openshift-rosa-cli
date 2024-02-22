package subnets

import (
	"strings"

	vpcClient "github.com/openshift-qe/openshift-rosa-cli/aws/vpc"
	awsV2 "github.com/openshift-qe/openshift-rosa-cli/pkg/aws_client/aws_v2"
	. "github.com/openshift-qe/openshift-rosa-cli/pkg/log"
	"github.com/spf13/cobra"
)

var args struct {
	region string
	zones  string
	vpcID  string
	tags   string
}

var Cmd = &cobra.Command{
	Use:   "subnets",
	Short: "Create subnets",
	Long:  "Create subnets.",
	Example: `  # Create a pair of subnets named "mysubne-"
  ocmqe create subnets --name-prefix=mysubnets --region us-east-2 --vpc-id <vpc id>`,

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
	Cmd.MarkFlagRequired("region")
	flags.StringVarP(
		&args.zones,
		"zones",
		"",
		"",
		"cidr of the vpc",
	)
	Cmd.MarkFlagRequired("zones")
	flags.StringVarP(
		&args.vpcID,
		"vpc-id",
		"",
		"",
		"Create a pair of subnets to",
	)
	Cmd.MarkFlagRequired("vpc-id")
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
	zones := strings.Split(args.zones, ",")
	for _, zone := range zones {
		subnetMap, err := vpc.PreparePairSubnetByZone(zone)
		if err != nil {
			panic(err)
		}
		for subnetType, subnet := range subnetMap {
			LogInfo("ZONE %s %s SUBNET: %s", zone, strings.ToUpper(subnetType), subnet.ID)
		}

	}

}
