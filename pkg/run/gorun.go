package run

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func GoModTidy(modroot string) (string, error) {
	log.Println("Running go mod tidy ...")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = modroot
	if bytes, err := cmd.CombinedOutput(); err != nil {
		return strings.TrimSpace(string(bytes)), err
	}
	return "", nil
}

func GoGetModule(name, version, modroot string) (string, error) {
	cmd := exec.Command("go", "get", fmt.Sprintf("%s@%s", name, version))
	cmd.Dir = modroot
	if bytes, err := cmd.CombinedOutput(); err != nil {
		return strings.TrimSpace(string(bytes)), err
	}
	return "", nil
}

func GoModEditReplaceModule(nameOld, nameNew, version, modroot string) (string, error) {
	cmd := exec.Command("go", "mod", "edit", "-dropreplace", fmt.Sprintf("%s", nameOld))
	cmd.Dir = modroot
	if bytes, err := cmd.CombinedOutput(); err != nil {
		return strings.TrimSpace(string(bytes)), fmt.Errorf("Error running go command to drop replace modules: %w", err)
	}

	cmd = exec.Command("go", "mod", "edit", "-replace", fmt.Sprintf("%s=%s@%s", nameOld, nameNew, version))
	cmd.Dir = modroot
	if bytes, err := cmd.CombinedOutput(); err != nil {
		return strings.TrimSpace(string(bytes)), fmt.Errorf("Error running go command to replace modules: %w", err)
	}
	return "", nil
}

func GoModEditRequireModule(name, version, modroot string) (string, error) {
	cmd := exec.Command("go", "mod", "edit", "-droprequire", fmt.Sprintf("%s", name))
	cmd.Dir = modroot
	if bytes, err := cmd.CombinedOutput(); err != nil {
		return strings.TrimSpace(string(bytes)), err
	}

	cmd = exec.Command("go", "mod", "edit", "-require", fmt.Sprintf("%s@%s", name, version))
	cmd.Dir = modroot
	if bytes, err := cmd.CombinedOutput(); err != nil {
		return strings.TrimSpace(string(bytes)), err
	}
	return "", nil
}
