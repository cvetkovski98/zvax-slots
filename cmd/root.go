package cmd

import (
	"fmt"

	"github.com/cvetkovski98/zvax-slots/internal/config"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Short: "Slots microservice",
	Long:  `Slots microservice`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Slots microservice")
	},
}

func init() {
	cobra.OnInitialize(configure)
	root.AddCommand(runCommand)
	root.AddCommand(seedCommand)
}

func configure() {
	if err := config.LoadConfig("config.dev.yaml"); err != nil {
		panic(err)
	}
}

func Execute() error {
	return root.Execute()
}
