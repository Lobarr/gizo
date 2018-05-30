package cli

import (
	"github.com/gizo-network/gizo/core"
	"github.com/spf13/cobra"
)

func init() {
	cleardbCmd.Flags().StringVarP(&env, "env", "e", "dev", "clear dev bc")
}

var cleardbCmd = &cobra.Command{
	Use:   "cleardb [flag]",
	Short: "Clears db",
	Run: func(cmd *cobra.Command, args []string) {
		if env == "dev" {
			core.RemoveDataPath()
		} else {
			core.RemoveDataPath()
		}
	},
}
