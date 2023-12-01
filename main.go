package main

import (
	"github.com/hectorj2f/godep-remedy/cmd"
	"github.com/hectorj2f/godep-remedy/pkg/constants"
)

func main() {
	cmd.Execute(cmd.Version{Version: constants.GoDepRemedyVersion, BuildDate: constants.GoDepRemedyBuildDate})
}
