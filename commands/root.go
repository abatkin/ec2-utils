package commands

import (
	"ec2-utils/aws"
	util2 "ec2-utils/display"
	"github.com/spf13/cobra"
)

func BuildRootCommand() *cobra.Command {

	rootCmd := &cobra.Command{
		Use:   "ec2",
		Short: "EC2 utilities",
	}

	awsOptions := aws.BuildAwsOptions(rootCmd)
	displayOptions := util2.BuildDisplayOptions(rootCmd)

	rootCmd.AddCommand(setupCommand(buildEc2ListCommand(), awsOptions, displayOptions))

	return rootCmd
}
