package cache

import (
	"archive/zip"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateFileWithDirectory(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp directory: %v", err)
		}
	}()

	filePath := filepath.Join(tempDir, "testdir", "testfile.txt")
	file, err := createFileWithDirectory(filePath)
	if err != nil {
		t.Errorf("createFileWithDirectory failed: %v", err)
	} else {
		if err := file.Close(); err != nil {
			t.Errorf("Failed to close file: %v", err)
		}
	}

	_, err = os.Stat(filePath)
	if err != nil {
		t.Errorf("File not created: %v", err)
	}
}

func TestLoadEtagCache(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp directory: %v", err)
		}
	}()

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
	err = os.WriteFile(etagFilePath, []byte("test-etag"), 0644)
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
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp directory: %v", err)
		}
	}()

	etagFilePath := filepath.Join(tempDir, "etag.txt")
	err = writeEtagCache(etagFilePath, "test-etag")
	if err != nil {
		t.Errorf("writeEtagCache failed: %v", err)
	}

	data, err := os.ReadFile(etagFilePath)
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
		if _, err := fmt.Fprint(w, "test content"); err != nil {
			t.Logf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp directory: %v", err)
		}
	}()

	// Test when no cache exists
	filePath, err := FetchFile(server.URL)
	if err != nil {
		t.Errorf("FetchFile failed when no cache exists: %v", err)
	}
	data, err := os.ReadFile(filePath)
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
	data, err = os.ReadFile(filePath)
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
		if _, err := fmt.Fprint(w, "new content"); err != nil {
			t.Logf("Failed to write response: %v", err)
		}
	})

	filePath, err = FetchFile(server.URL)
	if err != nil {
		t.Errorf("FetchFile failed when cache exists but etag doesn't match: %v", err)
	}
	data, err = os.ReadFile(filePath)
	if err != nil {
		t.Errorf("Failed to read cached file: %v", err)
	}
	if string(data) != "new content" {
		t.Errorf("Cached file content incorrect: %s", string(data))
	}
}

func TestExtractZipFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp directory: %v", err)
		}
	}()

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
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Failed to close file: %v", err)
		}
	}()

	zipWriter := zip.NewWriter(file)
	defer func() {
		if err := zipWriter.Close(); err != nil {
			log.Printf("Failed to close zip writer: %v", err)
		}
	}()

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
