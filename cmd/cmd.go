package cmd

import (
	"github.com/RiemaLabs/indexer-light/config"
	"github.com/spf13/cobra"
)

type Cmd struct {
}

func NewCmd() *Cmd {
	return &Cmd{}
}
func (c *Cmd) Init() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "Nubit Light Indexer",
		Short: "",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	rootCmd.Flags().IntVarP(&config.Config.StartHeight, "StartHeight", "s", config.Config.StartHeight, "Initial block height")
	rootCmd.Flags().StringVarP(&config.Config.StartBlockHash, "StartBlockHash", "", config.Config.StartBlockHash, "Initial block hash")
	rootCmd.Flags().IntVarP(&config.Config.MinimalCheckPoint, "MinimalCheckPoint", "", config.Config.MinimalCheckPoint, "Minimum amount of CheckPoint data required")
	rootCmd.Flags().StringVarP(&config.Config.BitCoinRpc.Host, "BitCoinRpc", "", config.Config.BitCoinRpc.Host, "BitCoin rpc endpoint")
	rootCmd.Flags().IntVarP(&config.Config.CommitteeIndexer.TimeOut, "TimeOut", "", config.Config.CommitteeIndexer.TimeOut, "CommitteeIndexer call timeout")
	return rootCmd
}
