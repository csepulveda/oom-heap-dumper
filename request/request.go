package request

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func RequestAndSave(pid, port int) (string, error) {
	url := fmt.Sprintf("http://localhost:%d/debug/pprof/heap", port)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to make request to %s: %v", url, err)
	}
	defer resp.Body.Close()

	// Check if the response status code is OK
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

	//get hostname
	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("failed to get hostname: %v", err)
	}

	//get date with hours and minutes
	date := time.Now().Format("2006-01-02-15-04")

	// Define the filename based on pid and port
	filename := fmt.Sprintf("%s-%d-%d-%s.heap", hostname, pid, port, date)
	filepath := filepath.Join(".", filename) // Save in the current directory

	// Create the file
	file, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %v", filename, err)
	}
	defer file.Close()

	// Write the response body to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write response to file %s: %v", filename, err)
	}

	// Return the full path of the saved file
	return filepath, nil
}

func DeleteFile(filepath string) error {
	err := os.Remove(filepath)
	if err != nil {
		return fmt.Errorf("failed to delete file %s: %v", filepath, err)
	}
	return nil
}
