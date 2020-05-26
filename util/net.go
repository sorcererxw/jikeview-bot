package util

import (
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
)

func downloadFile(url string) (string, error) {
	// Create the file
	filepath := os.TempDir() + "/" + uuid.New().String()
	out, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}
	return filepath, nil
}
