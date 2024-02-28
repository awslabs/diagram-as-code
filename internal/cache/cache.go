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
	"path"
)

func FetchFile(url string) (string, error) {
	hashedUrl := md5.New()
	io.WriteString(hashedUrl, url)
	cacheFilePath := fmt.Sprintf("%s/%x-%s", cacheDir, hashedUrl.Sum(nil), path.Base(url))

	// Check cached same URL resource
	if _, err := os.Stat(cacheFilePath); err != nil {

		resp, err := http.Get(url)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		out, err := os.Create(cacheFilePath)
		if err != nil {
			return "", err
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
	}
	return cacheFilePath, nil
}

func ExtractZipFile(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hashedContent := md5.New()
	if _, err := io.Copy(hashedContent, f); err != nil {
		return "", nil
	}
	cacheFilePath := fmt.Sprintf("%s/%x-%s", cacheDir, hashedContent.Sum(nil), path.Base(filePath))
	if _, err := os.Stat(cacheFilePath); err != nil {

		// Open a zip archive.
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

			err = os.MkdirAll(path.Dir(outputFilename), os.ModePerm)
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

const cacheDir = "./.cache"
