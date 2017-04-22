package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version of this tool",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(GetVersion())
	},
}

func GetVersion() string {
	return fmt.Sprintf("v%s\n", Version)
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

