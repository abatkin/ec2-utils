package commands

import (
	util2 "ec2-utils/util"
	"github.com/spf13/cobra"
)



func BuildRootCommand() *cobra.Command {

	rootCmd := &cobra.Command{
		Use:   "ec2",
		Short: "EC2 utilities",
	}

	awsOptions := util2.BuildAwsOptions(rootCmd)
	displayOptions := util2.BuildDisplayOptions(rootCmd)

	rootCmd.AddCommand(buildEc2ListCommand(awsOptions, displayOptions))

	return rootCmd
}


