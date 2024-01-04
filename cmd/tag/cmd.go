package tag

import (
	"strings"

	awsV2 "github.com/openshift-qe/openshift-rosa-cli/pkg/aws_client/aws_v2"
	"github.com/spf13/cobra"
)

var args struct {
	resourceID string
	tags       string
	region     string
}
var Cmd = &cobra.Command{
	Use:   "tag",
	Short: "tag a resource",
	Long:  "Tag a resource with the resource ID",
	Example: `  #Tag a vpc with vpc ID
  ocmqe tag --resource-id <vpc id> --tags tag1:tagv,tag2:tagv2`,

	Run: run,
}

func init() {
	flags := Cmd.Flags()
	flags.SortFlags = false
	flags.StringVarP(
		&args.resourceID,
		"resource-id",
		"",
		"",
		"resource ID",
	)
	flags.StringVarP(
		&args.region,
		"region",
		"",
		"",
		"region ID",
	)
	flags.StringVarP(
		&args.tags,
		"tags",
		"",
		"",
		"key of the tag",
	)
	Cmd.MarkFlagRequired("resource-id")
	Cmd.MarkFlagRequired("key")
	Cmd.MarkFlagRequired("region")
}
func run(cmd *cobra.Command, _ []string) {
	console, err := awsV2.CreateAWSV2Client("", args.region)
	if err != nil {
		panic(err)
	}
	splitedTags := map[string]string{}

	for _, tag := range strings.Split(args.tags, ",") {
		tagPair := strings.Split(tag, ":")
		if len(tagPair) < 2 {
			tagPair = append(tagPair, "")
		}
		splitedTags[tagPair[0]] = tagPair[1]
	}
	_, err = console.TagResource(args.resourceID, splitedTags)
	if err != nil {
		panic(err)
	}
}
