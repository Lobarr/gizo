package cli

import (
	"os"

	"github.com/gizo-network/gizo/helpers"
	"github.com/gizo-network/gizo/p2p"
	"github.com/kpango/glg"
	"github.com/spf13/cobra"
)

func init() {
	dispatcherCmd.Flags().IntVarP(&port, "port", "p", 9999, "port to run dispatcher on")
}

var dispatcherCmd = &cobra.Command{
	Use:   "dispatcher",
	Short: "Spin up a dispatcher node",
	Run: func(cmd *cobra.Command, args []string) {
		helpers.Banner()
		if os.Getenv("ENV") == "dev" {
			glg.Log("Core: using dev blockchain")
		}
		d := p2p.NewDispatcher(port)
		d.Start()
	},
}
