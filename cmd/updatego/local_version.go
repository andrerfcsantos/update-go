package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// VersionInfo contains information about a local go installation
type VersionInfo struct {
	Version      string
	OS           string
	Architecture string
}

// GoInPath checks if the go command exists in PATH
func GoInPath() bool {
	_, err := exec.LookPath("go")
	if err != nil {
		return false
	}
	return true
}

// GoVersion returns information about the local version. Returns an error if the go command is not in PATH.
func GoVersion() (VersionInfo, error) {
	vString, err := goVersionCmd()
	if err != nil {
		return VersionInfo{}, fmt.Errorf("error getting go version: %w", err)
	}

	versionParts := strings.Split(vString, " ")

	if len(versionParts) != 4 {
		return VersionInfo{}, fmt.Errorf("'go version' returned a non-expected output with %d components instead of the expected %d", len(versionParts), 4)
	}

	version := versionParts[2]

	osArch := strings.Split(versionParts[3], "/")
	if len(osArch) != 2 {
		return VersionInfo{}, fmt.Errorf("could not split the <os/arch> component of 'go version' output (%s) into os and arch %#v", versionParts[3], osArch)
	}

	goos, goarch := osArch[0], osArch[1]

	return VersionInfo{
		Version:      version,
		OS:           goos,
		Architecture: goarch,
	}, nil

}

// goVersionCmd performs the 'go version' command and returns its output
func goVersionCmd() (string, error) {

	cmd := exec.Command("go", "version")

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("couldn't get output of 'go version' command: %w", err)
	}

	return strings.TrimSpace(string(out)), nil
}
