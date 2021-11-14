package commands

import (
	"context"
	"ec2-utils/cmd/ec2/util"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/spf13/cobra"
	"log"
)

type ListOptions struct {
	instanceIds   []string
	filterOptions []string
}

func buildEc2ListCommand(awsOptions *util.AwsOptions) *cobra.Command {
	listOptions := &ListOptions{}

	return &cobra.Command{
		Use: "list ec2 instances",
		Run: func(cmd *cobra.Command, args []string) {
			listEc2(awsOptions, listOptions)
		},
	}
}

func listEc2(awsOptions *util.AwsOptions, listOptions *ListOptions) {
	client := util.BuildAwsClient(awsOptions)
	requestContext, _ := context.WithTimeout(context.Background(), awsOptions.Timeout)
	instancesPaginator := ec2.NewDescribeInstancesPaginator(client, &ec2.DescribeInstancesInput{})
	for instancesPaginator.HasMorePages() {
		page, err := instancesPaginator.NextPage(requestContext)
		if err != nil {
			log.Fatalf("error in pagination: %v", err)
		}
		for _, reservation := range page.Reservations {
			for _, instance := range reservation.Instances {
				log.Printf("instance: %s", *instance.InstanceId)
			}
		}
	}

}
