package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// GoVersionsURL is the online source from where to fecth information about go versions.
const GoVersionsURL = "https://golang.org/dl/?mode=json"

// RemoteVersions represents a list of information about go versions.
type RemoteVersions []RemoteVersion

// RemoteVersion represents information about a Go version.
// This information should come from an online source like https://golang.org/dl/?mode=json
type RemoteVersion struct {
	Version string       `json:"version"`
	Stable  bool         `json:"stable"`
	Files   []RemoteFile `json:"files"`
}

// GetFile returns the file information of an installer of a given kind, for this version of go, for the specified GOOS e GOARCH.
func (r *RemoteVersion) GetFile(goos, goarch, kind string) (RemoteFile, bool) {

	for _, f := range r.Files {
		if f.Os == goos && f.Arch == goarch && f.Kind == kind {
			return f, true
		}
	}
	return RemoteFile{}, false
}

// OSArch represents a tuple of an OS and Architecture
type OSArch struct {
	OS   string
	Arch string
}

// OSArchs returns the list of pairs GOOS/GOARCH this version of go supports
func (r *RemoteVersion) OSArchs() []OSArch {

	var res []OSArch
	auxMap := make(map[OSArch]bool)

	for _, f := range r.Files {
		if f.Kind != "source" {
			auxMap[OSArch{
				OS:   f.Os,
				Arch: f.Arch,
			}] = true
		}
	}

	for k := range auxMap {
		res = append(res, k)
	}

	return res
}

// RemoteFile represents information about a remote file for a specific version of Go.
type RemoteFile struct {
	Filename string `json:"filename"`
	Os       string `json:"os"`
	Arch     string `json:"arch"`
	Version  string `json:"version"`
	Sha256   string `json:"sha256"`
	Size     int    `json:"size"`
	Kind     string `json:"kind"`
}

// GetMostRecentVersion returns information about the most recent version of Go available.
func GetMostRecentVersion() (RemoteVersion, error) {
	vs, err := FetchRemoteVersions()
	if err != nil {
		return RemoteVersion{}, fmt.Errorf("couldn't fetch go versions: %w", err)
	}

	if len(vs) == 0 {
		return RemoteVersion{}, fmt.Errorf("feteched go versions sucessfully, but the returned array came empty")
	}

	return vs[0], nil
}

// FetchRemoteVersions returns information about the all versions of Go available. Use the function GetMostRecentVersion
// to get the most recent of them.
func FetchRemoteVersions() (RemoteVersions, error) {
	resp, err := http.Get(GoVersionsURL)
	if err != nil {
		return RemoteVersions{}, fmt.Errorf("making GET request for go versions: %w", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var vs RemoteVersions
	err = json.Unmarshal(body, &vs)
	if err != nil {
		return RemoteVersions{}, fmt.Errorf("could not unmarshal version json: %w", err)
	}

	return vs, nil
}
