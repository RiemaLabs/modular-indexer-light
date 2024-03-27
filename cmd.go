package main

import (
	"github.com/RiemaLabs/modular-indexer-light/log"
	"github.com/spf13/cobra"
)

type RuntimeArguments struct {
	EnableTest           bool
	TestBlockHeightLimit uint
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
			if arguments.EnableTest && arguments.TestBlockHeightLimit != 0 {
				log.Info("Use the test mode and limit the max blockheight %d to avoid catching up to the real latest block.\n", arguments.TestBlockHeightLimit)
			}
			Execution(arguments)
		},
	}

	rootCmd.Flags().BoolVarP(&arguments.EnableTest, "test", "t", false, "Enable this flag to hijack the blockheight to test the service")
	rootCmd.Flags().UintVarP(&arguments.TestBlockHeightLimit, "blockheight", "", 0, "When -test enabled, you can set TestBlockHeightLimit as a fixed value you want.")
	return rootCmd
}
