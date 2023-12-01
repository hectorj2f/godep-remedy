package types

type Module struct {
	Name    string
	Version string
	Replace bool
	Require bool
}

type CommonGoModError struct {
	Name   string
	Occurs func(goModTidyOutput string) bool
}
