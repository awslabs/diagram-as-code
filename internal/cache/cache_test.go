package cache

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateFileWithDirectory(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	filePath := filepath.Join(tempDir, "testdir", "testfile.txt")
	file, err := createFileWithDirectory(filePath)
	if err != nil {
		t.Errorf("createFileWithDirectory failed: %v", err)
	} else {
		file.Close()
	}

	_, err = os.Stat(filePath)
	if err != nil {
		t.Errorf("File not created: %v", err)
	}
}

func TestLoadEtagCache(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	etagFilePath := filepath.Join(tempDir, "etag.txt")

	// Test when file doesn't exist
	etag, err := loadEtagCache(etagFilePath)
	if err != nil {
		t.Errorf("loadEtagCache failed when file doesn't exist: %v", err)
	}
	if etag != "" {
		t.Errorf("loadEtagCache returned non-empty string when file doesn't exist")
	}

	// Test when file exists
	err = ioutil.WriteFile(etagFilePath, []byte("test-etag"), 0644)
	if err != nil {
		t.Fatalf("Failed to create etag file: %v", err)
	}
	etag, err = loadEtagCache(etagFilePath)
	if err != nil {
		t.Errorf("loadEtagCache failed when file exists: %v", err)
	}
	if etag != "test-etag" {
		t.Errorf("loadEtagCache returned incorrect etag value: %s", etag)
	}
}

func TestWriteEtagCache(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	etagFilePath := filepath.Join(tempDir, "etag.txt")
	err = writeEtagCache(etagFilePath, "test-etag")
	if err != nil {
		t.Errorf("writeEtagCache failed: %v", err)
	}

	data, err := ioutil.ReadFile(etagFilePath)
	if err != nil {
		t.Errorf("Failed to read etag file: %v", err)
	}
	if string(data) != "test-etag" {
		t.Errorf("Etag file content incorrect: %s", string(data))
	}
}

func TestFetchFile(t *testing.T) {
	// Test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		etag := r.Header.Get("If-None-Match")
		if etag == "test-etag" {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Header().Set("Etag", `"test-etag"`)
		fmt.Fprint(w, "test content")
	}))
	defer server.Close()

	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test when no cache exists
	filePath, err := FetchFile(server.URL)
	if err != nil {
		t.Errorf("FetchFile failed when no cache exists: %v", err)
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Errorf("Failed to read cached file: %v", err)
	}
	if string(data) != "test content" {
		t.Errorf("Cached file content incorrect: %s", string(data))
	}

	// Test when cache exists and etag matches
	filePath, err = FetchFile(server.URL)
	if err != nil {
		t.Errorf("FetchFile failed when cache exists and etag matches: %v", err)
	}
	data, err = ioutil.ReadFile(filePath)
	if err != nil {
		t.Errorf("Failed to read cached file: %v", err)
	}
	if string(data) != "test content" {
		t.Errorf("Cached file content incorrect: %s", string(data))
	}

	// Test when cache exists but etag doesn't match
	server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		etag := r.Header.Get("If-None-Match")
		if etag == "new-etag" {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Header().Set("Etag", `"new-etag"`)
		fmt.Fprint(w, "new content")
	})

	filePath, err = FetchFile(server.URL)
	if err != nil {
		t.Errorf("FetchFile failed when cache exists but etag doesn't match: %v", err)
	}
	data, err = ioutil.ReadFile(filePath)
	if err != nil {
		t.Errorf("Failed to read cached file: %v", err)
	}
	if string(data) != "new content" {
		t.Errorf("Cached file content incorrect: %s", string(data))
	}
}

func TestExtractZipFile(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	zipFilePath := filepath.Join(tempDir, "test.zip")
	err = createTestZipFile(zipFilePath)
	if err != nil {
		t.Fatalf("Failed to create test zip file: %v", err)
	}

	extractedPath, err := ExtractZipFile(zipFilePath)
	if err != nil {
		t.Errorf("ExtractZipFile failed: %v", err)
	}

	expectedFilePath := filepath.Join(extractedPath, "test.txt")
	_, err = os.Stat(expectedFilePath)
	if err != nil {
		t.Errorf("Extracted file not found: %v", err)
	}
}

func createTestZipFile(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	testFileData := []byte("test content")
	f, err := zipWriter.Create("test.txt")
	if err != nil {
		return err
	}
	_, err = f.Write(testFileData)
	if err != nil {
		return err
	}

	return nil
}
