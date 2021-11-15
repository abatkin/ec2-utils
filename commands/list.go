package commands

import (
	utilAws "ec2-utils/aws"
	utilDisplay "ec2-utils/display"
	"errors"
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
	command.Flags().StringSliceVar(&c.instanceIds, "instance-id", []string{}, "Instance `id`(s) to return")
	command.Flags().StringToStringVar(&c.filterOptions, "filter", map[string]string{}, "`key=value` values for filter")
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
	var filter *ec2.DescribeInstancesInput
	var err error
	if filter, err = c.buildFilter(); err != nil {
		return nil, err
	}

	instancesPaginator := ec2.NewDescribeInstancesPaginator(client, filter)
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

func (c *ListEc2InstancesCommand) buildFilter() (*ec2.DescribeInstancesInput, error) {
	request := &ec2.DescribeInstancesInput{}

	if len(c.instanceIds) > 0 {
		if len(c.filterOptions) > 0 {
			return nil, errors.New("only instance IDs OR fields maybe specified to filter")
		}

		request.InstanceIds = c.instanceIds
	}  else if len(c.filterOptions) > 0 {
		filters := make([]types.Filter, len(c.filterOptions))
		i := 0
		for key, value := range c.filterOptions {
			filters[i] = types.Filter{
				Name:   &key,
				Values: strings.Split(value, ","),
			}
			i++
		}

		request.Filters = filters
	}

	return request, nil
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
