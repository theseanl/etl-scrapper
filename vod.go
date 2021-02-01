package snuetl

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/oopsguy/m3u8/dl"
)

func downloadVod(client *http.Client, id int64, filename string) error {
	// Check for existing files
	if fileExists(filepath.Join(filename, "main.ts")) {
		log.Printf("Skipping download of an existing vod %v, id: %v\n", filename, id)
		return nil
	}
	log.Printf("Downloading vod %v, id: %v\n", filename, id)
	url, err := getM3u8UrlForEtlHostedVod(client, id)
	if err != nil {
		return err
	}
	return downloadVodFromM3u8Url(client, url, filename)
}

var reJwplayerScriptDecl *regexp.Regexp = regexp.MustCompile(`(?m)jwplayer\.key\s=[\s\S]*?file\s*:\s*['"]([^'"]*)['"][\s\S]*?jwplayer\(\)`)
var reNoConvertedVOD *regexp.Regexp = regexp.MustCompile(`(?i)There\sare\sno\sconverted\sVOD`)

func getM3u8UrlForEtlHostedVod(client *http.Client, id int64) (string, error) {
	// Query http://etl.snu.ac.kr/mod/vod/viewer.php?id=XXXXXX, get m3u8 URL.
	var match [][]byte = nil
	counter := 0
	for true {
		res, _ := client.Get(etlPath(fmt.Sprintf("/mod/vod/viewer.php?id=%v", id)))
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		res.Body.Close()
		match = reJwplayerScriptDecl.FindSubmatch(body)
		if match == nil {
			// Sometimes it provides "No converted VOD available" message.
			// Retry in 2 seconds until a converted VOD become available.
			if reNoConvertedVOD.Match(body) {
				// retry in 2 seconds.
				time.Sleep(2 * time.Second)
				counter++
				log.Printf("Retrying download: %v\n", counter)
				continue
			}
			return "", fmt.Errorf("Cannot get url for vod %v\n"+
				"==== Received http response below ==== %s", id, string(body))
		}
		break
	}
	return string(match[1]), nil
}

func downloadVodFromM3u8Url(client *http.Client, url string, filename string) error {
	os.MkdirAll(filepath.Dir(filename), os.ModePerm)
	downloader, err := dl.NewTask(filename, url)
	if err != nil {
		return err
	}
	if err := downloader.Start(25); err != nil {
		return err
	}
	return nil
}
