package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type Version struct {
	Version   string `json:"Version"`
	BuildDate string `json:"BuildDate"`
}

var out = os.Stdout
var errOut = os.Stderr
var in = os.Stdin

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "godep-remedy",
	Short: "godep-remedy cli",
	Args:  cobra.NoArgs,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version for secret",
}

func RootCmd() *cobra.Command {
	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version Version) {
	// set version string
	rootCmd.SetVersionTemplate(versionString(version))

	// also need to set Version to get cobra to print it
	rootCmd.Version = version.Version

	// add a version command
	versionCmd.Run = func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(out, versionString(version))
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func versionString(version Version) string {
	bytes, _ := json.MarshalIndent(version, "", "    ")
	return string(bytes) + "\n"
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.DisableAutoGenTag = true
}
