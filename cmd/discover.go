package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/hectorj2f/godep-remedy/pkg/discover"
)

type discoverCLIFlags struct {
	module  string
	modroot string
}

var discoverFlags discoverCLIFlags

var discoverCmd = &cobra.Command{
	SilenceUsage:  true,
	SilenceErrors: true,
	Use:           "discover [modroot] --module",
	Short:         "Discover go module versions",
	RunE: func(cmd *cobra.Command, args []string) error {
		if discoverFlags.module != "" {
			vers, err := discover.DiscoverModuleVersions(discoverFlags.modroot, discoverFlags.module)
			if err != nil {
				log.Println(err)
				return err
			}
			log.Printf("Found versions: %v\n", vers)

		} else {
			_, err := discover.DiscoverModules(discoverFlags.modroot, []string{})
			if err != nil {
				log.Println(err)
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(discoverCmd)

	flagSet := discoverCmd.Flags()
	flagSet.StringVar(&discoverFlags.module, "module", "", "go module path")
	flagSet.StringVar(&discoverFlags.modroot, "modroot", "", "go mod root path")
}
