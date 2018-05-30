package cli

import (
	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/helpers"
	"github.com/spf13/cobra"
)

func init() {
	cleardbCmd.Flags().StringVarP(&env, "env", "e", "dev", "clear dev bc")
}

var cleardbCmd = &cobra.Command{
	Use:   "cleardb [flag]",
	Short: "Clears db",
	Run: func(cmd *cobra.Command, args []string) {
		logger := helpers.Logger()
		if env == "dev" {
			logger.Info("Core: wiping dev blockchain")
			core.RemoveDataPath()
		} else {
			logger.Info("Core: wiping blockchain")
			core.RemoveDataPath()
		}
	},
}
