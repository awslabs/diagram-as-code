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
)

func FetchFile(url string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	hashedUrl := md5.New()
	io.WriteString(hashedUrl, url)
	cacheFilePath := filepath.Join(homeDir, ".cache", "awsdac", fmt.Sprintf("%x-%s", hashedUrl.Sum(nil), filepath.Base(url)))

	// Check cached same URL resource
	if _, err := os.Stat(cacheFilePath); err != nil {

		resp, err := http.Get(url)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 400 {
			return "", fmt.Errorf("Failed to fetch file %s: http status %d", url, resp.StatusCode)
		}

		err = os.MkdirAll(filepath.Dir(cacheFilePath), os.ModePerm)
		if err != nil {
			return "", err
		}

		out, err := os.Create(cacheFilePath)
		if err != nil {
			return "", err
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return "", err
		}
	}
	return cacheFilePath, nil
}

func ExtractZipFile(filePath string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hashedContent := md5.New()
	if _, err := io.Copy(hashedContent, f); err != nil {
		return "", nil
	}
	cacheFilePath := filepath.Join(homeDir, ".cache", "awsdac", fmt.Sprintf("%x-%s", hashedContent.Sum(nil), filepath.Base(filePath)))
	if _, err := os.Stat(cacheFilePath); err != nil {

		r, err := zip.OpenReader(filePath)
		if err != nil {
			return "", err
		}
		for _, f := range r.File {
			rc, err := f.Open()
			if err != nil {
				return "", err
			}

			outputFilename := fmt.Sprintf("%s/%s", cacheFilePath, f.Name)

			err = os.MkdirAll(filepath.Dir(outputFilename), os.ModePerm)
			if err != nil {
				return "", err
			}

			fo, err := os.Create(outputFilename)
			if err != nil {
				return "", err
			}
			_, err = io.Copy(fo, rc)
			if err != nil {
				return "", err
			}
			rc.Close()
			fo.Close()
		}

		defer r.Close()
	}

	return cacheFilePath, nil
}
