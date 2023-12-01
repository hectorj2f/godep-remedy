package cmd

import (
	"log"
	"os"
	"path"
	"strings"

	"github.com/hectorj2f/godep-remedy/pkg/constants"
	"github.com/hectorj2f/godep-remedy/pkg/types"
	"github.com/hectorj2f/godep-remedy/pkg/update"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

type updateCLIFlags struct {
	gomodules             string
	modroot               string
	tidy                  bool
	remedy                bool
	validateModuleVersion bool
}

type module struct {
	Name    string
	Version string
	Replace bool
	Require bool
}

var updateFlags updateCLIFlags

var updateCmd = &cobra.Command{
	SilenceUsage:  true,
	SilenceErrors: true,
	Use:           "update [gomodule@version]",
	Short:         "update a comma-separated list of go modules to update",
	RunE: func(cmd *cobra.Command, args []string) error {
		if updateFlags.gomodules == "" {
			log.Println("Usage: godep-remedy -modules=<module@version>,...")
			os.Exit(1)
		}
		packages := strings.Split(updateFlags.gomodules, ",")
		modules := map[string]*types.Module{}
		for _, pkg := range packages {
			parts := strings.Split(pkg, "@")
			if len(parts) != 2 {
				log.Println("Usage: godep-remedy -modules=<module@version>,...")
				os.Exit(1)
			}
			modules[parts[0]] = &types.Module{
				Name:    parts[0],
				Version: parts[1],
			}
		}
		if err := Run(updateFlags.modroot, modules, updateFlags.tidy, updateFlags.remedy, updateFlags.validateModuleVersion); err != nil {
			log.Printf("ERROR: run %v\n", err)
			os.Exit(1)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	flagSet := updateCmd.Flags()
	flagSet.StringVar(&updateFlags.gomodules, "modules", "", "comma separated list of go modules with versions")
	flagSet.StringVar(&updateFlags.modroot, "modroot", "", "go mod root path")
	flagSet.BoolVar(&updateFlags.tidy, "tidy", false, "go mod tidy flag")
	flagSet.BoolVar(&updateFlags.remedy, "remedy", true, "EXPERIMENTAL: enable auto-solve the go mod conflicts")
	flagSet.BoolVar(&updateFlags.validateModuleVersion, "validate-module-version", true, "check the module version and error if it founds a newer version")
}

func Run(modroot string, modules map[string]*types.Module, tidy, remedy, validateModuleVersion bool) error {
	modpath := path.Join(modroot, "go.mod")
	mod, err := update.ParseGoModfile(modpath)
	if err != nil {
		log.Printf("ERROR: unable to parse the go mod file with error: %v\n", err)
		return err
	}
	for _, replace := range mod.Replace {
		if replace != nil {
			if _, ok := modules[replace.New.Path]; ok {
				// pkg is already been replaced
				modules[replace.New.Path].Replace = true
				if semver.Compare(replace.New.Version, modules[replace.New.Path].Version) > 0 {
					log.Printf("WARNING: found replace package %s with version lower than the desired version %s", modules[replace.New.Path].Name, modules[replace.New.Path].Version)
				}
			}
		}
	}
	for _, require := range mod.Require {
		if require != nil {
			if _, ok := modules[require.Mod.Path]; ok {
				// pkg is already been required
				modules[require.Mod.Path].Require = true
				if semver.Compare(require.Mod.Version, modules[require.Mod.Path].Version) > 0 {
					log.Printf("WARNING: found require package %s with version lower than the desired version %s", require.Mod.Path, modules[require.Mod.Path].Version)
				}
			}
		}
	}
	// Attempt to fix the version
	for i := 1; i < constants.MaxRemedyAttempts; i++ {
		log.Println("Updating the modules...")
		if err := update.GoDepUpdate(modules, modroot, tidy, remedy, validateModuleVersion); err != nil {
			log.Printf("ERROR to go dep update: %v\n", err)
		}
		if err == nil {
			break
		}
	}
	return nil
}
