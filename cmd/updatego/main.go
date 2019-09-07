package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	mostRecent, err := GetMostRecentVersion()
	if err != nil {
		fmt.Printf("couldn't get most recent go version: %v", err)
		os.Exit(1)
	}

	if GoInPath() {
		fmt.Println("The go command was found in PATH. Checking for updates.")

		localVersion, err := GoVersion()
		if err != nil {
			fmt.Printf("could not get local go version: %v", err)
			os.Exit(1)
		}

		if mostRecent.Version == localVersion.Version {
			fmt.Printf("Your Go version (%v) is up to date.\n", localVersion.Version)
			os.Exit(0)
		} else {
			fmt.Printf("Your Go version (%v) is outdated. The most recent version is %v. Performing upgrade.\n", localVersion.Version, mostRecent.Version)
		}
	} else {
		fmt.Println("Go not found in PATH. Trying to install.")
	}

	err = PerformInstall(mostRecent)
	if err != nil {
		fmt.Printf("Error performing go instalation: %v", err)
		os.Exit(1)
	}

}

// PerformInstall executes the actions needed to install/upgrade the system for the specified target version.
func PerformInstall(targetVersion RemoteVersion) error {
	var err error
	var kind string
	var location string
	goos, goarch := runtime.GOOS, runtime.GOARCH

	switch {
	case goos == "windows" || goos == "darwin":
		kind = "installer"
	default:
		kind = "archive"
	}

	if f, ok := targetVersion.GetFile(goos, goarch, kind); ok {
		fmt.Printf("Downloading %s\n", f.Filename)
		location, err = DownloadFile(f.Filename, f.Sha256)
		if err != nil {
			return fmt.Errorf("downloading file: %w", err)
		}
	}

	switch kind {
	case "installer":
		err := RunInstaller(location)
		if err != nil {
			return fmt.Errorf("running installer: %w", err)
		}
	case "archive":
		err := ArchiveInstall(location)
		if err != nil {
			return fmt.Errorf("installing via tar.gz archive: %w", err)
		}
		fmt.Println("Go installed successfully!\nNote: Since the instalation was made via tarball, don't forget to make sure '/usr/local/go/bin' is in your path for the go command to work. For more information, check: https://golang.org/doc/install#tarball")
	default:
		return fmt.Errorf("kind of instalation '%s' unknown", kind)
	}

	fmt.Printf("Performing cleanup of downloaded file at '%s'\n", location)
	err = os.RemoveAll(location)
	if err != nil {
		fmt.Printf("could not remove file %s. You may need to remove it manually. Error: %v\n", location, err)
	}
	return nil
}

// RunInstaller runs a go installer on a given local path.
func RunInstaller(installerPath string) error {
	var args []string

	switch {
	case strings.HasSuffix(installerPath, ".msi"):
		args = []string{"MsiExec", "/package", installerPath}
	case strings.HasSuffix(installerPath, ".pkg"):
		args = []string{"sudo", "installer", "-pkg", installerPath, "-target", "/"}
	default:
		return fmt.Errorf("no suitable method was found for running the installer %s - expected an .msi or .pkg file", installerPath)
	}

	cmd := exec.Command(args[0], args[1:]...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("executing installer command: %w", err)
	}

	if !cmd.ProcessState.Success() {
		return fmt.Errorf("installer command exited with a non-sucess code (%d)", cmd.ProcessState.ExitCode())
	}
	return nil
}

// ArchiveInstall performs the tasks to install go from an tar.gz archive at the specified path
func ArchiveInstall(archivePath string) error {
	if !strings.HasSuffix(archivePath, ".tar.gz") {
		return fmt.Errorf("%s doesn't seem to be a tarball - expected a .tar.gz file for an archive installation", archivePath)
	}

	_, err := exec.LookPath("tar")
	if err != nil {
		return fmt.Errorf("could not locate the command 'tar' in PATH, which is required for instalation via archive")
	}

	if _, err := os.Stat("/usr/local/go"); !os.IsNotExist(err) {
		// '/usr/local/go' exists, remove old version
		err := os.RemoveAll("/usr/local/go")
		if err != nil {
			return fmt.Errorf("removing old version from /usr/local/go: %w", err)
		}
	}

	tarCmd := exec.Command("sudo", "tar", "-C", "/usr/local", "-xzf", archivePath)
	err = tarCmd.Run()
	if err != nil {
		return fmt.Errorf("error running 'tar' command: %w", err)
	}

	if !tarCmd.ProcessState.Success() {
		return fmt.Errorf("tar command exited with a non-sucess code (%d)", tarCmd.ProcessState.ExitCode())
	}

	return nil
}
