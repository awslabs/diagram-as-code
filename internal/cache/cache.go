// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cache

import (
	"archive/zip"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func FetchFile(url string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Cannot get home directory: %v", err)
	}

	hashedUrl := md5.New()
	io.WriteString(hashedUrl, url)
	cacheFilePath := filepath.Join(homeDir, ".cache", "awsdac", fmt.Sprintf("%x-%s", hashedUrl.Sum(nil), filepath.Base(url)))

	// Check cached same URL resource
	if _, err := os.Stat(cacheFilePath); err != nil {

		resp, err := http.Get(url)
		if err != nil {
			return "", fmt.Errorf("Cannot get HTTP resource(%s): %v", url, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 400 {
			return "", fmt.Errorf("Failed to fetch file %s: http status %d", url, resp.StatusCode)
		}

		err = os.MkdirAll(filepath.Dir(cacheFilePath), os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("Cannot create directory(%s): %v", filepath.Dir(cacheFilePath), err)
		}

		out, err := os.Create(cacheFilePath)
		if err != nil {
			return "", fmt.Errorf("Cannot create file(%s): %v", cacheFilePath, err)
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return "", fmt.Errorf("Cannot copy: %v", err)
		}
	}
	return cacheFilePath, nil
}

func ExtractZipFile(filePath string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Cannot get home directory: %v", err)
	}

	f, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("Cannot open file(%s): %v", filePath, err)
	}
	defer f.Close()

	hashedContent := md5.New()
	if _, err := io.Copy(hashedContent, f); err != nil {
		return "", fmt.Errorf("Cannot create md5 hash: %v", err)
	}
	cacheFilePath := filepath.Join(homeDir, ".cache", "awsdac", fmt.Sprintf("%x-%s", hashedContent.Sum(nil), filepath.Base(filePath)))
	if _, err := os.Stat(cacheFilePath); err != nil {

		r, err := zip.OpenReader(filePath)
		if err != nil {
			return "", fmt.Errorf("Cannot open file(%s): %v", filePath, err)
		}
		for _, f := range r.File {
			if strings.HasSuffix(f.Name, "/") {
				continue
			}
			rc, err := f.Open()
			if err != nil {
				return "", fmt.Errorf("Cannot open file(%s): %v", f.Name, err)
			}

			outputFilename := fmt.Sprintf("%s/%s", cacheFilePath, f.Name)

			err = os.MkdirAll(filepath.Dir(outputFilename), os.ModePerm)
			if err != nil {
				return "", fmt.Errorf("Cannot create directory(%s): %v", filepath.Dir(outputFilename), err)
			}

			fo, err := os.Create(outputFilename)
			if err != nil {
				return "", fmt.Errorf("Cannot create file(%s): %v", outputFilename, err)
			}
			_, err = io.Copy(fo, rc)
			if err != nil {
				return "", fmt.Errorf("Cannot copy: %v", err)
			}
			rc.Close()
			fo.Close()
		}

		defer r.Close()
	}

	return cacheFilePath, nil
}
