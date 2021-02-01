package snuetl

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func downloadAssignment(client *http.Client, id int64, instanceName string) error {
	log.Printf("Downloading Assignment %v, id: %v\n", instanceName, id)
	// http://etl.snu.ac.kr/mod/assign/view.php?id=XXXX
	resp, err := client.Get(etlPath(fmt.Sprintf("/mod/assign/view.php?id=%v", id)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	errors := []error{}
	doc.Find(`[id^="assign_files_tree"] [yuiconfig] a`).Each(func(i int, content *goquery.Selection) {
		href, ok := content.Attr("href")
		if !ok {
			log.Printf("Cannot find assignment attachment url for %v\n", id)
			return
		}
		filename := filepath.Join(instanceName, content.Text())
		filename = strings.TrimSuffix(filename, filepath.Ext(filename))

		err := downloadFile(client, href, filename)
		if err != nil {
			errors = append(errors, err)
		}
	})

	if len(errors) > 0 {
		return fmt.Errorf("%v", errors)
	}
	return nil
}
