package commands

import (
	utilAws "ec2-utils/aws"
	"ec2-utils/display"
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

func (c *ListEc2InstancesCommand) runCommand(awsOptions *utilAws.Options) ([]display.Item, error) {
	client := awsOptions.BuildAwsClient()
	requestContext := awsOptions.BuildRequestContext()
	var filter *ec2.DescribeInstancesInput
	var err error
	if filter, err = c.buildFilter(); err != nil {
		return nil, err
	}

	instancesPaginator := ec2.NewDescribeInstancesPaginator(client, filter)
	items := make([]display.Item, 0)
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

func (c *ListEc2InstancesCommand) defaultFields() []display.FieldInfo {
	return []display.FieldInfo{
		{Name: "id", Expression: "InstanceId"},
		{Name: "Name", Expression: "Tags:Name"},
		{Name: "State", Expression: "State.Name"},
	}
}

func (c *ListEc2InstancesCommand) headerName(fieldName display.FieldInfo) string {
	return fieldName.Name
}

func (c *ListEc2InstancesCommand) buildFilter() (*ec2.DescribeInstancesInput, error) {
	request := &ec2.DescribeInstancesInput{}

	if len(c.instanceIds) > 0 {
		if len(c.filterOptions) > 0 {
			return nil, errors.New("only instance IDs OR fields maybe specified to filter")
		}

		request.InstanceIds = c.instanceIds
	} else if len(c.filterOptions) > 0 {
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

func newEc2Item(instance types.Instance) display.Item {
	var tags = map[string]string{}
	for _, tag := range instance.Tags {
		tags[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
	}

	return &Ec2Item{
		Instance: instance,
		tags:     tags,
	}
}

func (e *Ec2Item) GetValue(fieldInfo display.FieldInfo) string {
	expression := fieldInfo.Expression
	if strings.Index(expression, "Tags:") == 0 {
		return e.tags[expression[5:]]
	}

	return display.ExtractFromExpression(expression, e)
}
