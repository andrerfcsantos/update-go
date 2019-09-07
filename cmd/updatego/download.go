package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// BaseDownloadURL represents the base path for go downloads
const BaseDownloadURL = "https://dl.google.com/go/"

// DownloadFile downloads a file with a given filename from the go downloads page. After the download also compares
// the checksum of the file with the one provided.
func DownloadFile(filename, expectedChecksum string) (string, error) {
	location := filepath.Join(".tmp", filename)

	err := downloadFile(location, BaseDownloadURL+filename)
	if err != nil {
		return location, fmt.Errorf("downloading file: %w", err)
	}

	f, err := os.Open(location)
	if err != nil {
		return location, fmt.Errorf("could not open downloaded file for checksum verification: %w", err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return location, fmt.Errorf("error copying data to hasher: %w", err)
	}

	actualChecksum := fmt.Sprintf("%x", h.Sum(nil))
	if expectedChecksum != actualChecksum {
		return location, fmt.Errorf("checksum verification failed (expected: %s | got: %s)", expectedChecksum, actualChecksum)
	}

	return location, nil
}

// downloadFile is the low-level function that actually performs the download of a file from a url to a given path
func downloadFile(path string, url string) error {

	d := filepath.Dir(path)

	err := os.MkdirAll(d, os.FileMode(0775))
	if err != nil {
		return fmt.Errorf("making directory to download go version (%s): %w", d, err)
	}

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error on GET request for file '%s': %w", filepath.Base(path), err)
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating file '%s': %w", filepath.Base(path), err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("error copying data from request to file '%s': %w", filepath.Base(path), err)
	}
	return nil
}
