package update

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hectorj2f/godep-remedy/pkg/discover"
	"github.com/hectorj2f/godep-remedy/pkg/types"
)

// remediate attempts to resolve `go mod tidy` issues such as modules that do not contain expected packages.
func remediate(modules map[string]*types.Module, modroot string, tidy bool, p types.CommonGoModError, output string) error {
	if p.Name == "ModuleDoesNotContainPackage" {
		re := regexp.MustCompile(`: module (.*?) found \(v(.*?)\), but does not contain package`)
		match := re.FindStringSubmatch(output)
		log.Printf("Detected error with module %s and version %s\n", match[1], match[2])
		moduleName := strings.Replace(match[1], "@latest", "", -1)
		moduleVer := match[2]

		verList, err := discover.DiscoverModuleVersions(modroot, moduleName)
		if err != nil {
			log.Println(err)
		}

		if _, ok := modules[moduleName]; ok {
			modules[moduleName].Replace = true
			modules[moduleName].Require = false
			idx := 0
			for it, v := range verList {
				if v == modules[moduleName].Version {
					idx = it
					break
				}
			}
			if idx > 0 {
				modules[moduleName].Version = verList[idx-1]
			}
			log.Printf("Module changed %v to version: %v", modules[moduleName], verList[idx-1])

		} else {
			idx := 0
			for it, v := range verList {
				if v == moduleVer {
					idx = it
					break
				}
			}
			if idx > 0 {
				modules[moduleName].Version = verList[idx-1]
			}
			modules[moduleName] = &types.Module{
				Replace: true,
				Version: moduleVer,
			}

			log.Printf("Module changed %v to version: %v", modules[moduleName], moduleVer)
		}

		// restore go mod files
		//if err := restoreGoModfile(modroot, dname); err != nil {
		//	return err
		//}
		return fmt.Errorf("Force restart with different version")
	}
	return nil
}
