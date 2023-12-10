package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/hectorj2f/godep-remedy/pkg/run"
	"github.com/hectorj2f/godep-remedy/pkg/types"
	"github.com/spf13/cobra"
)

type replaceCLIFlags struct {
	gomodule string
	modroot  string
	tidy     bool
}

var replaceFlags replaceCLIFlags

var replaceCmd = &cobra.Command{
	SilenceUsage:  true,
	SilenceErrors: true,
	Use:           "replace [module-to-replace=newgomodule@version]",
	Short:         "use to replace a single go module by another go module",
	RunE: func(cmd *cobra.Command, args []string) error {
		if replaceFlags.gomodule == "" {
			log.Println("Usage: godep-remedy -modules=<module@version>,...")
			os.Exit(1)
		}
		parts := strings.Split(replaceFlags.gomodule, "=")
		if len(parts) != 2 {
			log.Println("Usage: godep-remedy -module-to-replace=<module=newgomodule@version>,...")
			os.Exit(1)
		}
		moduOld := &types.Module{
			Name: parts[0],
		}
		newModParts := strings.Split(parts[1], "@")

		if len(newModParts) != 2 {
			log.Println("Usage: godep-remedy -module-to-replace=<module=newgomodule@version>,...")
			os.Exit(1)
		}
		moduNew := &types.Module{
			Name:    newModParts[0],
			Version: newModParts[1],
		}
		log.Println("Running go mod edit replace ...")
		if _, err := run.GoModEditReplaceModule(moduOld.Name, moduNew.Name, moduNew.Version, replaceFlags.modroot); err != nil {
			log.Printf("ERROR: go mod replace %v\n", err)
			os.Exit(1)
		}

		if replaceFlags.tidy {
			if _, err := run.GoModTidy(replaceFlags.modroot); err != nil {
				log.Printf("ERROR: go mod tidy %v\n", err)
				os.Exit(1)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(replaceCmd)

	flagSet := replaceCmd.Flags()
	flagSet.StringVar(&replaceFlags.gomodule, "module-to-replace", "", "a go module to replace and its new go module with the version")
	flagSet.StringVar(&replaceFlags.modroot, "modroot", "", "go mod root path")
	flagSet.BoolVar(&replaceFlags.tidy, "tidy", false, "go mod tidy flag")
}
