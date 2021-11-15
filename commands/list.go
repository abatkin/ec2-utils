package commands

import (
	util2 "ec2-utils/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

type ListOptions struct {
	instanceIds   []string
	filterOptions map[string]string
}

func buildEc2ListCommand(awsOptions *util2.AwsOptions, displayOptions *util2.DisplayOptions) *cobra.Command {
	listOptions := &ListOptions{}

	return &cobra.Command{
		Use: "list ec2 instances",
		Run: func(cmd *cobra.Command, args []string) {
			listEc2(awsOptions, listOptions, displayOptions)
		},
	}
}

type Ec2Item struct {
	types.Instance
	tags map[string]string
}

func (e *Ec2Item) GetValue(name string) string {
	switch {
	case strings.Index(name, "tags.") == 0:
		tagName := name[5:]
		return e.tags[tagName]
	case name == "id":
		return aws.ToString(e.InstanceId)
	case name == "state":
		return string(e.State.Name)
	default:
		return ""
	}
}

func listEc2(awsOptions *util2.AwsOptions, listOptions *ListOptions, displayOptions *util2.DisplayOptions) {
	client := awsOptions.BuildAwsClient()
	requestContext := awsOptions.BuildRequestContext()

	instancesPaginator := ec2.NewDescribeInstancesPaginator(client, &ec2.DescribeInstancesInput{})
	items := make([]util2.Item, 0)
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

	displayOptions.Render(getFields(displayOptions), items)
}

func newEc2Item(instance types.Instance) util2.Item {
	var tags = map[string]string{}
	for _, tag := range instance.Tags {
		tags[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
	}

	return &Ec2Item{
		Instance: instance,
		tags:     tags,
	}
}

var DefaultFields = []util2.Field{
	{FieldName: "id", Heading: "InstanceId"},
	{FieldName: "tags.Name", Heading: "Name"},
	{FieldName: "state", Heading: "State"},
}

func getFields(displayOptions *util2.DisplayOptions) []util2.Field {
	if len(displayOptions.Fields) == 0 {
		return DefaultFields
	}

	fields := make([]util2.Field, len(displayOptions.Fields))
	for i, fieldName := range displayOptions.Fields {
		fields[i] = util2.Field{FieldName: fieldName, Heading: buildHeading(fieldName)}
	}
	return fields
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
