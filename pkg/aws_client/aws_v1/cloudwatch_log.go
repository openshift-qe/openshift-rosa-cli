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
		fmt.Printf("Got error describe log group:%s ", err)
	}
	return *output, err
}

func (client *AWSClient) DescribeLogStreamByName(logGroupName string) (cloudwatchlogs.DescribeLogStreamsOutput, error) {
	output, err := client.cloudWatchLogsClient.DescribeLogStreams(&cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: &logGroupName,
	})
	if err != nil {
		fmt.Printf("Got error describe log stream: %s", err)
	}
	return *output, err
}

func (client *AWSClient) DeleteLogGroupByName(logGroupName string) (cloudwatchlogs.DeleteLogGroupOutput, error) {
	output, err := client.cloudWatchLogsClient.DeleteLogGroup(&cloudwatchlogs.DeleteLogGroupInput{
		LogGroupName: &logGroupName,
	})
	if err != nil {
		fmt.Printf("Got error delete log group: %s", err)
	}
	return *output, err
}
