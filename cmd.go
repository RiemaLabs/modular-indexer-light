package main

import (
	"github.com/RiemaLabs/modular-indexer-light/log"
	"github.com/spf13/cobra"
)

type RuntimeArguments struct {
	EnableTest     bool
	EnableDAReport bool
}

func NewRuntimeArguments() *RuntimeArguments {
	return &RuntimeArguments{}
}

func (arguments *RuntimeArguments) MakeCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "Nubit Light Indexer",
		Short: "Activates the Nubit Light Indexer with optional services.",
		Long: `
Light Indexer is an essential component of the Nubit Modular Indexer architecture.
It enables typical users to verify Bitcoin meta-protocols without requiring substantial computing resources.
This command offers multiple flags to tailor the indexer's functionality according to the user's needs.
		`,
		Run: func(cmd *cobra.Command, args []string) {
			if arguments.EnableTest {
				log.Info("Use the test mode to setup gin as debug mode.")
			}

			if arguments.EnableDAReport {
				log.Info("Report to DA is enabled.")
			} else {
				log.Info("Report to DA is disabled.")
			}

			Execution(arguments)
		},
	}

	rootCmd.Flags().BoolVarP(&arguments.EnableTest, "test", "t", false, "Enable this flag to hijack the blockheight to test the service")
	rootCmd.Flags().BoolVarP(&arguments.EnableDAReport, "report", "", true, "Enable this flag to upload verified checkpoint to DA")
	return rootCmd
}
