package aws_v1

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

func (client *AWSClient) DescribeLogGroupsByName(logGroupName string) (cloudwatchlogs.DescribeLogGroupsOutput, error) {
	output, err := client.cloudWatchLogsClient.DescribeLogGroups(&cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: &logGroupName,
	})
	if err != nil {
		fmt.Println("Got error describe log group: ", err)
	}
	return *output, err
}

func (client *AWSClient) DescribeLogStreamByName(logGroupName string) (cloudwatchlogs.DescribeLogStreamsOutput, error) {
	output, err := client.cloudWatchLogsClient.DescribeLogStreams(&cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: &logGroupName,
	})
	if err != nil {
		fmt.Println("Got error describe log stream: ", err)
	}
	return *output, err
}
