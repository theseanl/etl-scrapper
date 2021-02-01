package snuetl

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/PuerkitoBio/goquery"
)

func downloadURLLink(client *http.Client, id int64, instancename string) error {
	// Check for existing files
	if fileExists(instancename + ".url") {
		log.Printf("Skipping url shortcut for id %v\n", id)
		return nil
	}
	log.Printf("Creating shortcut for id %v\n", id)

	// Visit http://etl.snu.ac.kr/mod/url/view.php?id=XXXXX
	req, err := http.NewRequest("GET", etlPath(fmt.Sprintf("/mod/url/view.php?id=%v", id)), nil)
	if err != nil {
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	doc, _ := goquery.NewDocumentFromReader(res.Body)

	anchor := doc.Find("#maincontent ~ .urlworkaround > a").First()

	href, _ := anchor.Attr("href")
	os.MkdirAll(filepath.Dir(instancename), os.ModePerm)
	out, err := os.Create(instancename + ".url")
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = out.WriteString(fmt.Sprintf(`[InternetShortcut]
URL=%s
`, href))

	return err
}
