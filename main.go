package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	url := ""
	path := "./"
	urlPtr := flag.String("url", "", "The url to download a zip file from")
	pathPtr := flag.String("path", "./", "The path where files should be extracted")
	intervalPtr := flag.String("interval", "", "The interval for how often the file should be downloaded, 15m, 6h, 1d, etc.")

	flag.Parse()

	if urlPtr == nil || *urlPtr == "" {
		log.Fatal("Url needs to be provided")
	} else {
		url = *urlPtr
	}
	if pathPtr != nil {
		path = *pathPtr
	}

	downloadAndUnzip(url, path)

	if intervalPtr == nil || *intervalPtr == "" {
		return
	}
	interval, err := time.ParseDuration(*intervalPtr)
	if err != nil {
		log.Fatal(err)
	}
	ticker := time.NewTicker(interval)
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			downloadAndUnzip(url, path)
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func downloadAndUnzip(url string, path string) {
	// https://stackoverflow.com/a/50539327
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		log.Fatal(err)
	}

	// Read all the files from zip archive
	for _, zipFile := range zipReader.File {
		log.Println("Reading file:", zipFile.Name)
		err := readZipFile(zipFile, path)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func readZipFile(zf *zip.File, dest string) error {
	log.Println(dest)
	err := os.MkdirAll(dest, 0755)
	if err != nil {
		return err
	}

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			err := os.MkdirAll(path, f.Mode())
			if err != nil {
				return err
			}
		} else {
			err := os.MkdirAll(filepath.Dir(path), f.Mode())
			if err != nil {
				return err
			}
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	err = extractAndWriteFile(zf)
	if err != nil {
		return err
	}

	return nil
}
