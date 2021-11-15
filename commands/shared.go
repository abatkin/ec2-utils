package commands

import (
	utilAws "ec2-utils/aws"
	"ec2-utils/display"
	"github.com/spf13/cobra"
	"log"
)

type Ec2Command interface {
	optionSetup(command *cobra.Command)
	usage() string
	description() string
	runCommand(awsOptions *utilAws.Options) ([]display.Item, error)
	defaultFieldNames() []string
	headerName(field string) string
}

func setupCommand(command Ec2Command, awsOptions *utilAws.Options, displayOptions *display.Options) *cobra.Command {
	cobraCommand := &cobra.Command{
		Use: command.usage(),
		Short: command.description(),
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			if result, err := command.runCommand(awsOptions); err != nil {
				log.Fatalf("error running command: %v", err)
			} else {
				var fields []string
				if len(displayOptions.Fields) > 0 {
					fields = displayOptions.Fields
				} else {
					fields = command.defaultFieldNames()
				}
				displayOptions.Render(fields, command.headerName, result)
			}
		},
	}
	command.optionSetup(cobraCommand)

	return cobraCommand
}
