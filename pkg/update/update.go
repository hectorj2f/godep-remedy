package update

import (
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/hectorj2f/godep-remedy/pkg"
	"github.com/hectorj2f/godep-remedy/pkg/discover"
	"github.com/hectorj2f/godep-remedy/pkg/run"
	"github.com/hectorj2f/godep-remedy/pkg/types"
	"golang.org/x/mod/modfile"
)

func revertGoModChanges(modpath string, content []byte) error {
	f, err := os.OpenFile(modpath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	if _, err := f.Write(content); err != nil {
		f.Close() // ignore error; Write error takes precedence
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

func copyGoModFile(src, dst string) error {
	_, err := exec.Command("cp", "-r", src, dst).Output()
	if err != nil {
		return err
	}
	return nil
}

func backupGoModfile(modroot string) (string, error) {
	modpath := path.Join(modroot, "go.mod")
	dname, err := os.MkdirTemp("", "dep-remediate")
	if err != nil {
		return dname, err
	}
	if err := copyGoModFile(modpath, dname); err != nil {
		return "", err
	}
	if err := copyGoModFile(path.Join(modroot, "go.sum"), dname); err != nil {
		return "", err
	}

	return dname, nil
}

func restoreGoModfile(modroot, dname string) error {
	if err := copyGoModFile(path.Join(dname, "go.mod"), modroot); err != nil {
		return err
	}
	if err := copyGoModFile(path.Join(dname, "go.sum"), modroot); err != nil {
		return err
	}

	return nil
}

func GoDepUpdate(modules map[string]*types.Module, modroot string, tidy, remedy, validateModuleVersion bool) error {
	dname, err := backupGoModfile(modroot)
	if err != nil {
		return err
	}

	for _, v := range modules {
		log.Printf("Remediate with module: %v and version %s\n", v.Name, v.Version)
		if validateModuleVersion {
			isLatest, err := discover.IsLatestModuleVersions(modroot, v.Name, v.Version)
			if err != nil {
				return err
			}
			if !isLatest {
				log.Println(err)
				os.Exit(1)
			}
		}
	}

	if output, err := update(modules, modroot, tidy); err != nil {
		cmps := pkg.FindCommonGoModErrors(output)
		log.Printf("Found common go mod errors on remediate: %v", cmps)

		// restore go mod files
		//if err := restoreGoModfile(modroot, dname); err != nil {
		//	return err
		//}
	}

	if tidy {
		if output, err := run.GoModTidy(modroot); err != nil {
			cmps := pkg.FindCommonGoModErrors(output)
			log.Printf("Found common go mod errors: %v", cmps)
			for _, p := range cmps {
				remediate(modules, modroot, tidy, p, output)
			}
			// restore go mod files
			if err := restoreGoModfile(modroot, dname); err != nil {
				return err
			}
		}

	}

	return nil
}

func update(modules map[string]*types.Module, modroot string, tidy bool) (string, error) {
	for k, pkg := range modules {
		log.Printf("Get package: %s\n", k)
		if pkg.Replace {
			log.Println("Running go mod edit replace ...")
			if output, err := run.GoModEditReplaceModule(pkg.Name, pkg.Name, pkg.Version, modroot); err != nil {
				log.Println(err)
				return output, err
			}
		} else if pkg.Require {
			log.Println("Running go mod edit require ...")
			if output, err := run.GoModEditRequireModule(pkg.Name, pkg.Version, modroot); err != nil {
				log.Println(err)
				return output, err
			}
		} else {
			log.Println("Running go get ...")
			if output, err := run.GoGetModule(pkg.Name, pkg.Version, modroot); err != nil {
				log.Println(err)
				return output, err
			}
		}
	}

	return "", nil
}

func ParseGoModfile(path string) (*modfile.File, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	mod, err := modfile.Parse("go.mod", content, nil)
	if err != nil {
		return nil, err
	}

	return mod, nil
}
