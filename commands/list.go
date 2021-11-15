package commands

import (
	utilAws "ec2-utils/aws"
	utilDisplay "ec2-utils/display"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/spf13/cobra"
	"strings"
)

type ListEc2InstancesCommand struct {
	instanceIds   []string
	filterOptions map[string]string
}

func (c *ListEc2InstancesCommand) optionSetup(command *cobra.Command) {
	// TODO implement this
}

func (c *ListEc2InstancesCommand) usage() string {
	return "list [--filter name=value ... | --instance-id id ...]"
}

func (c *ListEc2InstancesCommand) description() string {
	return "List EC2 Instances"
}

func (c *ListEc2InstancesCommand) runCommand(awsOptions *utilAws.Options) ([]utilDisplay.Item, error) {
	client := awsOptions.BuildAwsClient()
	requestContext := awsOptions.BuildRequestContext()

	instancesPaginator := ec2.NewDescribeInstancesPaginator(client, &ec2.DescribeInstancesInput{})
	items := make([]utilDisplay.Item, 0)
	for instancesPaginator.HasMorePages() {
		page, err := instancesPaginator.NextPage(requestContext)
		if err != nil {
			return nil, err
		}
		for _, reservation := range page.Reservations {
			for _, instance := range reservation.Instances {
				items = append(items, newEc2Item(instance))
			}
		}
	}
	return items, nil
}

func (c *ListEc2InstancesCommand) defaultFieldNames() []string {
	return []string{"id", "tags.Name", "state"}
}

func (c *ListEc2InstancesCommand) headerName(fieldName string) string {
	if strings.Index(fieldName, "Tags.") == 0 {
		return fieldName[5:]
	} else if fieldName == "id" {
		return "InstanceId"
	} else if len(fieldName) >= 2 {
		return strings.ToUpper(fieldName[0:1]) + fieldName[1:]
	} else {
		return fieldName
	}
}

func buildEc2ListCommand() Ec2Command {
	listOptions := &ListEc2InstancesCommand{}
	return listOptions
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
