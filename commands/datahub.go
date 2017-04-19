package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:	"datahub",
	Short:	"datahub is the CLI for DataHub.io",
	Long:	"",
	Run:	func(cmd *cobra.Command, args []string) {
		fmt.Println("DataHub 🐘  a home for your data packages ❒ ❒ ❒ ")
	},
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

