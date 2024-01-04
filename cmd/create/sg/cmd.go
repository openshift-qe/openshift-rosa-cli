package sg

import (
	"fmt"
	"strings"

	vpcClient "github.com/openshift-qe/openshift-rosa-cli/aws/vpc"
	awsV2 "github.com/openshift-qe/openshift-rosa-cli/pkg/aws_client/aws_v2"
	. "github.com/openshift-qe/openshift-rosa-cli/pkg/log"
	"github.com/spf13/cobra"
)

var args struct {
	region     string
	count      int
	vpcID      string
	tags       string
	namePrefix string
}

var Cmd = &cobra.Command{
	Use:   "security-groups",
	Short: "Create security-groups",
	Long:  "Create security-groups.",
	Example: `# Create a number of security groups"
  ocmqe create security-groups --name-prefix=mysg --region us-east-2 --vpc-id <vpc id>`,

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
		"Region of the security groups",
	)
	flags.StringVarP(
		&args.namePrefix,
		"name-prefix",
		"",
		"",
		"Name prefix of the security groups, they will be named with <prefix>-0,<prefix>-1",
	)

	flags.IntVarP(
		&args.count,
		"count",
		"",
		0,
		"Additional security  groups  number going to be created to the vpc",
	)
	flags.StringVarP(
		&args.vpcID,
		"vpc-id",
		"",
		"",
		"vpc id going to create the addtional security groups to",
	)
	Cmd.MarkFlagRequired("vpc-id")
	Cmd.MarkFlagRequired("region")
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
	preparedSGs := []string{}
	createdsgNum := 0
	sgDescription := "This security group is created for OCM testing"
	protocol := "tcp"
	for createdsgNum < args.count {
		sgName := fmt.Sprintf("%s-%d", args.namePrefix, createdsgNum)
		sg, err := vpc.AWSClient.CreateSecurityGroup(vpc.VpcID, sgName, sgDescription)
		if err != nil {
			panic(err)
		}
		groupID := *sg.GroupId
		cidrPortsMap := map[string]int32{
			vpc.CIDRValue: 8080,
			"0.0.0.0/0":   22,
		}
		for cidr, port := range cidrPortsMap {
			_, err = vpc.AWSClient.AuthorizeSecurityGroupIngress(groupID, cidr, protocol, port, port)
			if err != nil {
				panic(err)
			}
		}

		preparedSGs = append(preparedSGs, *sg.GroupId)
		createdsgNum++
	}
	LogInfo("ADDITIONAL SECURITY GROUPS: %s", strings.Join(preparedSGs, ","))
}
