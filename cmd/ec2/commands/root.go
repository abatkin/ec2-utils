package commands

import (
	"ec2-utils/cmd/ec2/util"
	"github.com/spf13/cobra"
)



func BuildRootCommand() *cobra.Command {

	rootCmd := &cobra.Command{
		Use:   "ec2",
		Short: "EC2 utilities",
	}

	awsOptions := util.BuildAwsOptions(rootCmd)
	displayOptions := util.BuildDisplayOptions(rootCmd)

	rootCmd.AddCommand(buildEc2ListCommand(awsOptions, displayOptions))

	return rootCmd
}


