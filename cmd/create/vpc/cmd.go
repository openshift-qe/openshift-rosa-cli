package vpc

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	vpcClient "gitlab.cee.redhat.com/openshift-group-I/ocm_aws/aws/vpc"
	. "gitlab.cee.redhat.com/openshift-group-I/ocm_aws/pkg/log"
)

var args struct {
	region       string
	name         string
	cidr         string
	tags         string
	findExisting bool
}
var Cmd = &cobra.Command{
	Use:   "vpc",
	Short: "Create vpc",
	Long:  "Create vpc.",
	Example: `  # Create a vpc named "myvpc"
  ocmqe create vpc --name=myvpc --region us-east-2`,

	Run: run,
}

func init() {
	flags := Cmd.Flags()
	flags.SortFlags = false
	flags.StringVarP(
		&args.name,
		"name",
		"n",
		"",
		"Name of the vpc",
	)
	flags.StringVarP(
		&args.region,
		"region",
		"",
		"",
		"Region of the vpc",
	)
	flags.StringVarP(
		&args.cidr,
		"cidr",
		"",
		"",
		"cidr of the vpc",
	)
	flags.StringVarP(
		&args.tags,
		"tags",
		"",
		"",
		"tags of the vpc, fmt tagName:tagValue,tagName2:tagValue2",
	)
	flags.BoolVarP(
		&args.findExisting,
		"find-existing",
		"",
		false,
		"Find the vpc with same name from current region. if not exsiting, create a new one",
	)
}
func run(cmd *cobra.Command, _ []string) {
	vpc, err := vpcClient.PrepareVPC(args.name, args.region, args.cidr, args.findExisting)
	if err != nil {
		LogError(err.Error())
		os.Exit(1)
	}

	LogInfo("VPC ID: %s", vpc.VpcID)
	LogInfo("VPC REGION: %s", vpc.Region)
	LogInfo("VPC NAME: %s", vpc.VPCName)
	var tagMap map[string]string
	if args.tags != "" {
		_tags := strings.Split(args.tags, ",")
		tagMap = map[string]string{}
		for _, tag := range _tags {
			splited := strings.Split(tag, ":")

			if len(splited) == 2 {
				tagMap[splited[0]] = tagMap[splited[1]]
			} else {
				tagMap[splited[0]] = ""
			}

		}
		vpc.AWSClient.TagResource(vpc.VpcID, tagMap)
	}
}
