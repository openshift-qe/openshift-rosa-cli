package tag

import (
	awsV2 "github.com/openshift-qe/openshift-rosa-cli/pkg/aws_client/aws_v2"
	"github.com/spf13/cobra"
)

var args struct {
	region     string
	resourceID string
	tagKey     string
	tagValue   string
}
var Cmd = &cobra.Command{
	Use:   "tag",
	Short: "Delete tag",
	Long:  "Delete tag.",
	Example: `  # Delete a tag from the resource
  ocmqe delete tag --resource-id <vpc id> --region us-east-2 --tag-key key --tag-value value`,

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
		"Region of the resource",
	)
	flags.StringVarP(
		&args.resourceID,
		"resource-id",
		"",
		"",
		"id of the resource",
	)
	flags.StringVarP(
		&args.tagKey,
		"tag-key",
		"",
		"",
		"tag key of the resource",
	)
	flags.StringVarP(
		&args.tagValue,
		"tag-value",
		"",
		"",
		"tag value of the resource",
	)

	Cmd.MarkFlagRequired("resource-id")
	Cmd.MarkFlagRequired("region")
	Cmd.MarkFlagRequired("tag-key")
	Cmd.MarkFlagRequired("tag-value")
}
func run(cmd *cobra.Command, _ []string) {
	console, err := awsV2.CreateAWSV2Client("", args.region)
	if err != nil {
		panic(err)
	}
	_, err = console.RemoveResourceTag(args.resourceID, args.tagKey, args.tagValue)
	if err != nil {
		panic(err)
	}
}
