package commands

import (
	utilAws "ec2-utils/aws"
	"ec2-utils/display"
	"github.com/spf13/cobra"
	"log"
)

type optionSetupFunction func(command *cobra.Command)
type commandFunction func(awsOptions *utilAws.Options) ([]display.Item, error)

type ec2Command struct {
	optionSetup optionSetupFunction
	usage string
	description string
	command commandFunction
	defaultFieldNames []string
	header display.HeaderFunction
}

func setupCommand(command ec2Command, awsOptions *utilAws.Options, displayOptions *display.Options) *cobra.Command {
	cobraCommand := &cobra.Command{
		Use: command.usage,
		Short: command.description,
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			if result, err := command.command(awsOptions); err != nil {
				log.Fatalf("error running command: %v", err)
			} else {
				var fields []string
				if len(displayOptions.Fields) > 0 {
					fields = displayOptions.Fields
				} else {
					fields = command.defaultFieldNames
				}
				displayOptions.Render(fields, command.header, result)
			}
		},
	}
	command.optionSetup(cobraCommand)

	return cobraCommand
}
