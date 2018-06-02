package cli

import (
	"github.com/gizo-network/gizo/helpers"
	"github.com/spf13/cobra"
)

var gizoCmd = &cobra.Command{
	Use:     "gizo [command]",
	Short:   "Job scheduling system build using blockchain",
	Long:    `Decentralized distributed system built using blockchain technology to provide a marketplace for users to trade their processing power in reward for ethereum`,
	Args:    cobra.MinimumNArgs(1),
	Version: "1.0.0",
}

//Execute boostraps all commands
func Execute() {
	gizoCmd.AddCommand(workerCmd, dispatcherCmd, cleardbCmd)
	if err := gizoCmd.Execute(); err != nil {
		helpers.Logger().Fatal(err)
	}
}
