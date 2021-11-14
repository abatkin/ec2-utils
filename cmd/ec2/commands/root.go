package commands

import (
	"ec2-utils/cmd/ec2/util"
	"github.com/spf13/cobra"
)



func BuildRootCommand() *cobra.Command {
	awsOptions := &util.AwsOptions{}

	rootCmd := cobra.Command{
		Use:   "ec2",
		Short: "EC2 utilities",
	}
	rootCmd.PersistentFlags().StringVar(&awsOptions.Profile, "profile", "", "AWS Profile")
	rootCmd.PersistentFlags().StringVar(&awsOptions.Region, "region", "", "AWS Region")
	rootCmd.PersistentFlags().DurationVar(&awsOptions.Timeout, "timeout", 0, "Timeout")
	rootCmd.PersistentFlags().StringVar(&awsOptions.Proxy, "proxy", "", "Proxy (set to 'none' to force-disable)")

	rootCmd.AddCommand(buildEc2ListCommand(awsOptions))

	return &rootCmd
}

