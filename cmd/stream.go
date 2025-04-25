package cmd

import "github.com/spf13/cobra"

var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: "Stream DB into another DB with different options",
	Long:  `This command streams the contents of this DB into another DB with the given options.`,
	RunE:  stream,
}

func stream(cmd *cobra.Command, args []string) error {

	return nil
}
