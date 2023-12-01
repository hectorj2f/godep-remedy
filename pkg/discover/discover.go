package discover

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"

	"github.com/Masterminds/semver/v3"
	sem "golang.org/x/mod/semver"
)

type Module struct {
	name string
	from *semver.Version
	to   *semver.Version
}

type moduleVersions struct {
	Path      string   `json:Path`
	Versions  []string `json:Versions`
	GoVersion string   `json:GoVersion`
}

func shouldIgnore(name, from, to string, ignoreNames []string) bool {
	for _, ig := range ignoreNames {
		if strings.Contains(name, ig) {
			log.Printf("Ignore module::: name: %s - from: %s - to: %s", name, from, to)
			return true
		}
	}
	return false
}
func IsLatestModuleVersions(modroot, name, version string) (bool, error) {
	args := []string{
		"list",
		"-m",
		"-mod=readonly",
		"-versions",
		"-json",
		name,
	}

	cmd := exec.Command("go", args...)
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("Error running go list to discover module %s versions: %s \n %w", name, string(bytes), err)
	}
	m := moduleVersions{}
	err = json.Unmarshal(bytes, &m)
	if err != nil {
		return false, err
	}

	for _, v := range m.Versions {
		if !strings.Contains(v, "-rc") && !strings.Contains(v, "-RC") && !strings.Contains(v, "-alpha") && !strings.Contains(v, "-beta") {
			if sem.Compare(v, version) > 0 {
				log.Printf("Update version for module: %s isn't the latest %s\n", name, v)
				return false, nil
			}
		}
	}
	return true, nil
}

func DiscoverModuleVersions(modroot, name string) ([]string, error) {
	args := []string{
		"list",
		"-m",
		"-mod=readonly",
		"-versions",
		"-json",
		name,
	}

	cmd := exec.Command("go", args...)
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		return []string{}, fmt.Errorf("Error running go list to discover module versions: %s \n %w", string(bytes), err)
	}
	m := moduleVersions{}
	// log.Printf("Content discover module versions \n%s\n", string(bytes))
	err = json.Unmarshal(bytes, &m)
	if err != nil {
		return []string{}, err
	}

	filteredVersions := make([]string, 0)
	for it, v := range m.Versions {
		// FIXME: semver only stable version
		if !strings.Contains(v, "-rc") && !strings.Contains(v, "-RC") && !strings.Contains(v, "-alpha") && !strings.Contains(v, "-beta") {
			filteredVersions = append(filteredVersions, m.Versions[it])
		}
	}
	return filteredVersions, nil
}

func DiscoverModules(modroot string, ignoreNames []string) (map[string]Module, error) {
	args := []string{
		"list",
		"-u",
		"-mod=readonly",
		"-f",
		"'{{if (and (not (or .Main .Indirect)) .Update)}}{{.Path}}: {{.Version}} -> {{.Update.Version}}{{end}}'",
		"-m",
		"all",
	}

	cmd := exec.Command("go", args...)
	//cmd.Env = os.Environ()
	// Disable Go workspace mode, otherwise this can cause trouble
	// See issue https://github.com/oligot/go-mod-upgrade/issues/35
	//cmd.Env = append(cmd.Env, "GOWORK=off")
	cmd.Dir = modroot

	bytes, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("Error running go list to discover modules: %s \n %w", string(bytes), err)
	}

	if err != nil {
		return nil, fmt.Errorf("Error running go command to discover modules: %w", err)
	}

	split := strings.Split(string(bytes), "\n")
	modules := map[string]Module{}
	re := regexp.MustCompile(`'(.+): (.+) -> (.+)'`)
	for _, x := range split {
		if x != "''" && x != "" {
			matched := re.FindStringSubmatch(x)
			if len(matched) < 4 {
				return nil, fmt.Errorf("Couldn't parse module %s", x)
			}
			name, from, to := matched[1], matched[2], matched[3]
			log.Printf("Found module::: name: %s - from: %s - to: %s\n", name, from, to)

			if shouldIgnore(name, from, to, ignoreNames) {
				continue
			}
			fromversion, err := semver.NewVersion(from)
			if err != nil {
				return nil, err
			}
			toversion, err := semver.NewVersion(to)
			if err != nil {
				return nil, err
			}
			d := Module{
				name: name,
				from: fromversion,
				to:   toversion,
			}
			modules[name] = d
		}
	}
	return modules, nil
}
