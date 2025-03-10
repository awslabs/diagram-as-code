// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cache

import (
	"archive/zip"
	"crypto/md5"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func createFileWithDirectory(filePath string) (*os.File, error) {
	err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("cannot create directory(%s): %v", filepath.Dir(filePath), err)
	}

	out, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot create file(%s): %v", filePath, err)
	}
	return out, nil
}

func writeFile(outputFilename string, fi *zip.File) error {
	rc, err := fi.Open()
	if err != nil {
		return fmt.Errorf("cannot open file: %v", err)
	}
	defer rc.Close()


	fo, err := createFileWithDirectory(outputFilename)
	if err != nil {
		return fmt.Errorf("cannot create file with directory: %v", err)
	}
	defer fo.Close()

	_, err = io.Copy(fo, rc)
	if err != nil {
		return fmt.Errorf("cannot copy: %v", err)
	}
	return nil
}

func loadEtagCache(etagFilePath string) (string, error) {
	// Check Etag file
	if _, err := os.Stat(etagFilePath); err == nil {
		f, err := os.Open(etagFilePath)
		if err != nil {
			return "", fmt.Errorf("cannot open Etag file(%s): %v", etagFilePath, err)
		}
		defer f.Close()

		bytes, err := io.ReadAll(f)
		if err != nil {
			return "", fmt.Errorf("cannot read Etag file(%s): %v", etagFilePath, err)
		}
		return string(bytes), nil
	}
	return "", nil
}

func writeEtagCache(etagFilePath, etag_value string) error {
	out, err := createFileWithDirectory(etagFilePath)
	if err != nil {
		return fmt.Errorf("cannot create file with directory: %v", err)
	}
	defer out.Close()

	d := []byte(etag_value)
	_, err = out.Write(d)
	if err != nil {
		return fmt.Errorf("cannot write Etag file(%s): %v", etagFilePath, err)
	}
	return nil
}

func FetchFile(url string) (string, error) {
	log.Infof("[internal/cache/cache.go] FetchFile %s", url)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot get home directory: %v", err)
	}

	hashedUrl := md5.New()
	io.WriteString(hashedUrl, url)

	etagFilePath := filepath.Join(homeDir, ".cache", "awsdac", "etag", fmt.Sprintf("%x-%s", hashedUrl.Sum(nil), filepath.Base(url)))
	cacheFilePath := filepath.Join(homeDir, ".cache", "awsdac", fmt.Sprintf("%x-%s", hashedUrl.Sum(nil), filepath.Base(url)))

	cached_etag_value := ""
	// Check cached same URL resource
	if _, err := os.Stat(cacheFilePath); err == nil {
		cached_etag_value, err = loadEtagCache(etagFilePath)
		if err != nil {
			return "", fmt.Errorf("cannot load Etag Cache: %v", err)
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("cannot create HTTP request(%s): %v", url, err)
	}

	if cached_etag_value != "" {
		log.Infof("[internal/cache/cache.go] Found previous Etag cache. Use HTTP Etag value %s", cached_etag_value)
		req.Header.Add("If-None-Match", cached_etag_value)
	}

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("cannot get HTTP resource(%s): %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return "", fmt.Errorf("failed to fetch file %s: http status %d", url, resp.StatusCode)
	}

	etag_value := ""
	if etag, ok := resp.Header["Etag"]; ok && len(etag) > 0 {
		z := strings.SplitN(etag[0], "/", 2)
		etag_value = z[0]
		if z[0] == "W" { // weak validator was returned.
			etag_value = z[1]
		}
	}
	log.Infof("[internal/cache/cache.go] Server respond with HTTP %d", resp.StatusCode)

	if resp.StatusCode == 302 && cached_etag_value == "" {
		return "", fmt.Errorf("remote server is responding with an HTTP 304 even though no If-none-match header was added to the request")
	}

	if resp.StatusCode == 302 && cached_etag_value != etag_value {
		return "", fmt.Errorf("remote server is responding with an HTTP 304 even though mismatch between Etag response header and If-none-Match request header")
	}

	// save remote resource to local if no local cache or etag mismatch or server doesn't send etag
	if cached_etag_value == "" || etag_value == "" || cached_etag_value != etag_value {
		out, err := createFileWithDirectory(cacheFilePath)
		if err != nil {
			return "", fmt.Errorf("cannot create file with directory: %v", err)
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return "", fmt.Errorf("cannot copy: %v", err)
		}

		// save as Etag
		if etag_value != "" {
			log.Infof("[internal/cache/cache.go] Server respond with Etag. Save Etag value %s", etag_value)
			err := writeEtagCache(etagFilePath, etag_value)
			if err != nil {
				return "", fmt.Errorf("cannot write Etag cache(%s): %v", etagFilePath, err)
			}
		}
	} else {
		log.Infof("[internal/cache/cache.go] Use cache based on matched HTTP Etag")
	}
	return cacheFilePath, nil
}

func ExtractZipFile(filePath string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot get home directory: %v", err)
	}

	f, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("cannot open file(%s): %v", filePath, err)
	}
	defer f.Close()

	hashedContent := md5.New()
	if _, err := io.Copy(hashedContent, f); err != nil {
		return "", fmt.Errorf("cannot create md5 hash: %v", err)
	}
	cacheFilePath := filepath.Join(homeDir, ".cache", "awsdac", fmt.Sprintf("%x-%s", hashedContent.Sum(nil), filepath.Base(filePath)))
	if _, err := os.Stat(cacheFilePath); err != nil {

		r, err := zip.OpenReader(filePath)
		if err != nil {
			return "", fmt.Errorf("cannot open file(%s): %v", filePath, err)
		}
		for _, f := range r.File {
			if strings.HasSuffix(f.Name, "/") {
				continue
			}
			outputFilename := fmt.Sprintf("%s/%s", cacheFilePath, f.Name)

			err := writeFile(outputFilename, f)
			if err != nil {
				return "", fmt.Errorf("cannot write file(%s): %v", f.Name, err)
			}
		}

		defer r.Close()
	}

	return cacheFilePath, nil
}
