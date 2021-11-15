package commands

import (
	utilAws "ec2-utils/aws"
	utilDisplay "ec2-utils/display"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

type Options struct {
	instanceIds   []string
	filterOptions map[string]string
}

func buildEc2ListCommand() ec2Command {
	listOptions := &Options{}

	return ec2Command{
		optionSetup: func(command *cobra.Command) {
			// TODO: add cli options using listOptions
		},
		usage:             "list [--filter name=value ... | --instance-id id ...]",
		description:       "List EC2 Instances",
		command: func(awsOptions *utilAws.Options) ([]utilDisplay.Item, error) {
			return listEc2(awsOptions, listOptions)
		},
		defaultFieldNames: []string{"id", "tags.Name", "state"},
		header:            buildHeading,
	}
}

type Ec2Item struct {
	types.Instance
	tags map[string]string
}

func newEc2Item(instance types.Instance) utilDisplay.Item {
	var tags = map[string]string{}
	for _, tag := range instance.Tags {
		tags[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
	}

	return &Ec2Item{
		Instance: instance,
		tags:     tags,
	}
}

func (e *Ec2Item) GetValue(name string) string {
	switch {
	case strings.Index(name, "tags.") == 0:
		return e.tags[name[5:]]
	case name == "id":
		return aws.ToString(e.InstanceId)
	case name == "state":
		return string(e.State.Name)
	default:
		return ""
	}
}

func listEc2(awsOptions *utilAws.Options, listOptions *Options) ([]utilDisplay.Item, error) {
	client := awsOptions.BuildAwsClient()
	requestContext := awsOptions.BuildRequestContext()

	instancesPaginator := ec2.NewDescribeInstancesPaginator(client, &ec2.DescribeInstancesInput{})
	items := make([]utilDisplay.Item, 0)
	for instancesPaginator.HasMorePages() {
		page, err := instancesPaginator.NextPage(requestContext)
		if err != nil {
			log.Fatalf("error in pagination: %v", err)
		}
		for _, reservation := range page.Reservations {
			for _, instance := range reservation.Instances {
				items = append(items, newEc2Item(instance))
			}
		}
	}
	return items, nil
}

func buildHeading(name string) string {
	if strings.Index(name, "Tags.") == 0 {
		return name[5:]
	} else if name == "id" {
		return "InstanceId"
	} else if len(name) >= 2 {
		return strings.ToUpper(name[0:1]) + name[1:]
	} else {
		return name
	}
}
