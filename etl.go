package snuetl

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const etlHost = "http://etl.snu.ac.kr"

func etlPath(pathname string) string {
	return etlHost + pathname
}

func getCourseURL(courseid int) string {
	return etlPath(fmt.Sprintf("/course/view.php?id=%v", courseid))
}

func DownloadEtlContents(client *http.Client, courseid int) {
	res, err := client.Get(getCourseURL(courseid))
	if err != nil {
		fmt.Printf(`Cannot open course page for id: %v\n, error: %v`, courseid, err)
		return
	}
	defer res.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(res.Body)

	var count = 0
	doc.Find(".total_sections .section.main").Each(func(i int, section *goquery.Selection) {
		sectionName := textContent(section.Find(".sectionname"))
		section.Find("li.activity").Each(func(i int, activity *goquery.Selection) {
			attrs, _ := activity.Attr("class")
			idstr, _ := activity.Attr("id")
			if !strings.HasPrefix(idstr, "module-") {
				cantParseActivityID(idstr, i)
				return
			}
			id, err := strconv.ParseInt(idstr[7:], 10, 32)
			if err != nil {
				cantParseActivityID(idstr, i)
				return
			}

			instanceName := textContent(activity.Find(".instancename"))
			fileName := filepath.Join(sectionName, instanceName)
			classes := strings.Split(attrs, " ")

			var downloadError error
			if contains(classes, "ubfile") {
				downloadError = ubfile(client, id, fileName)
			} else if contains(classes, "vod") {
				downloadError = downloadVod(client, id, fileName)
			} else if contains(classes, "url") {
				downloadError = downloadURLLink(client, id, fileName)
			} else if contains(classes, "assign") {
				downloadError = downloadAssignment(client, id, fileName)
			} else {
				downloadError = fmt.Errorf(`Cannot determine content type for %v, got classes %v`, id, classes)
			}
			if downloadError != nil {
				fmt.Println(downloadError)
			}
			count++
		})
	})
	fmt.Printf("Processed %v sections.\n", count)
}

func cantParseActivityID(idstr string, index int) {
	fmt.Printf("Can't parse id for an activity ID from %v, index %v\n", idstr, index)
}

func textContent(s *goquery.Selection) string {
	var out = ""
	s.Contents().Not("*").Each(func(i int, s *goquery.Selection) {
		out += s.Text()
	})
	return out
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
