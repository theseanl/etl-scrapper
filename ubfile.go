package snuetl

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

func ubfile(client *http.Client, id int64, instanceName string) error {
	// http://etl.snu.ac.kr/mod/ubfile/view.php?id=XXXXX
	log.Printf("Downloading file %v, id: %v\n", instanceName, id)
	return downloadFile(client, etlPath(fmt.Sprintf("/mod/ubfile/view.php?id=%v", id)), instanceName)
}

func downloadFile(client *http.Client, url string, name string) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Use Content-Disposition header to determine file extension
	_, params, _ := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
	filename, ok := params["filename"]
	if !ok {
		filename = ""
	}

	filename = name + filepath.Ext(filename)

	// Check for existing files
	if fileExists(filename) {
		log.Printf("Skipping download of an existing file\n")
		return nil
	}

	os.MkdirAll(filepath.Dir(filename), os.ModePerm)
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
