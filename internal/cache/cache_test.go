package cache

import (
	"archive/zip"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestFetchFile(t *testing.T) {
	// Set up a test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test content"))
	}))
	defer server.Close()

	// Test case 1: Fetch file and cache it
	url := server.URL
	cacheFilePath, err := FetchFile(url)
	if err != nil {
		t.Errorf("FetchFile(%s) returned error: %v", url, err)
	}

	// Check if the cached file exists
	if _, err := os.Stat(cacheFilePath); err != nil {
		t.Errorf("Cached file not found: %s", cacheFilePath)
	}

	// Check the cached file content
	cachedContent, err := ioutil.ReadFile(cacheFilePath)
	if err != nil {
		t.Errorf("Failed to read cached file: %v", err)
	}
	if string(cachedContent) != "test content" {
		t.Errorf("Cached file content is incorrect: %s", cachedContent)
	}

	// Test case 2: Fetch file from cache
	cacheFilePathAgain, err := FetchFile(url)
	if err != nil {
		t.Errorf("FetchFile(%s) returned error: %v", url, err)
	}
	if cacheFilePathAgain != cacheFilePath {
		t.Errorf("Cached file path mismatch: %s != %s", cacheFilePathAgain, cacheFilePath)
	}
}

func TestExtractZipFile(t *testing.T) {
	// Create a temporary zip file
	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Errorf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	zipFilePath := filepath.Join(tempDir, "test.zip")
	err = createZipFile(zipFilePath, map[string]string{
		"file1.txt":     "content1",
		"dir/file2.txt": "content2",
	})
	if err != nil {
		t.Errorf("Failed to create zip file: %v", err)
	}

	// Test case 1: Extract zip file
	extractedDirPath, err := ExtractZipFile(zipFilePath)
	if err != nil {
		t.Errorf("ExtractZipFile(%s) returned error: %v", zipFilePath, err)
	}

	// Check if the extracted files exist
	file1Path := filepath.Join(extractedDirPath, "file1.txt")
	if _, err := os.Stat(file1Path); err != nil {
		t.Errorf("Extracted file not found: %s", file1Path)
	}

	file2Path := filepath.Join(extractedDirPath, "dir/file2.txt")
	if _, err := os.Stat(file2Path); err != nil {
		t.Errorf("Extracted file not found: %s", file2Path)
	}

	// Check the extracted file content
	extractedContent1, err := ioutil.ReadFile(file1Path)
	if err != nil {
		t.Errorf("Failed to read extracted file: %v", err)
	}
	if string(extractedContent1) != "content1" {
		t.Errorf("Extracted file content is incorrect: %s", extractedContent1)
	}

	extractedContent2, err := ioutil.ReadFile(file2Path)
	if err != nil {
		t.Errorf("Failed to read extracted file: %v", err)
	}
	if string(extractedContent2) != "content2" {
		t.Errorf("Extracted file content is incorrect: %s", extractedContent2)
	}

	// Test case 2: Extract zip file from cache
	extractedDirPathAgain, err := ExtractZipFile(zipFilePath)
	if err != nil {
		t.Errorf("ExtractZipFile(%s) returned error: %v", zipFilePath, err)
	}
	if extractedDirPathAgain != extractedDirPath {
		t.Errorf("Extracted directory path mismatch: %s != %s", extractedDirPathAgain, extractedDirPath)
	}
}

func createZipFile(zipFilePath string, files map[string]string) error {
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for path, content := range files {
		f, err := zipWriter.Create(path)
		if err != nil {
			return err
		}
		_, err = f.Write([]byte(content))
		if err != nil {
			return err
		}
	}

	return nil
}
