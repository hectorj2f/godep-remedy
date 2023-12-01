package pkg

import (
	"regexp"

	"github.com/hectorj2f/godep-remedy/pkg/types"
)

var CommonGoModErrors []types.CommonGoModError = []types.CommonGoModError{
	{
		Name: "ModuleDoesNotContainPackage",
		Occurs: func(goInstallOutput string) bool {
			r := regexp.MustCompile("module .* found .*, but does not contain package .*")
			return r.MatchString(goInstallOutput)
		},
	},

	{
		Name: "RequiredConflictError",
		Occurs: func(goInstallOutput string) bool {
			r := regexp.MustCompile("(?s)module declares its path as: .*\\n\\s+but was required as: .*")
			return r.MatchString(goInstallOutput)
		},
	},

	{
		Name: "ReplaceDirectivesError",
		Occurs: func(goInstallOutput string) bool {
			r := regexp.MustCompile("(?s)The go.mod file .* contains .* replace directives.")
			return r.MatchString(goInstallOutput)
		},
	},
}

func FindCommonGoModErrors(goInstallOutput string) []types.CommonGoModError {
	var problems []types.CommonGoModError

	for _, p := range CommonGoModErrors {
		if p.Occurs(goInstallOutput) {
			problems = append(problems, p)
		}
	}

	return problems
}
