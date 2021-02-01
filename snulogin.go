package snuetl

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
)

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.104 Safari/537.36"

// SnuLogin creates a session logged in into snu.ac.kr with the provided user credential.
func SnuLogin(username string, password string) (*http.Client, error) {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar:     jar,
		Timeout: time.Second * 1000,
	}

	req, _ := http.NewRequest("GET", "https://sso.snu.ac.kr/snu/ssologin.jsp", nil)
	req.Header.Set("User-Agent", userAgent)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	actionURL := findActionURL(string(body))

	postData := url.Values{}
	postData.Set("userid", username)
	postData.Set("password", password)
	postData.Set("si_redirect_address", "")
	postData.Set("lang_id", "ko")
	postData.Set("si_realm", "SnuUser1")
	postData.Set("id_save", "on")

	req, _ = http.NewRequest("POST", actionURL, strings.NewReader(postData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)

	res, err = client.Do(req)

	body, _ = ioutil.ReadAll(res.Body)
	bodyStr, _ := decodeToEUCKR(string(body))

	html, _ := html.Parse(strings.NewReader(bodyStr))
	secondPostData := url.Values{}
	findInputTags(html, func(key string, value string) {
		secondPostData.Set(key, value)
	})

	req, _ = http.NewRequest("POST", "https://sso.snu.ac.kr/nls3/fcs", strings.NewReader(secondPostData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)

	res, err = client.Do(req)
	res.Body.Close()
	if err != nil {
		return nil, err
	}
	return client, nil
}

var reSubmitURL = regexp.MustCompile(`f.action=["']([^"']*)["']`)

func findActionURL(text string) string {
	match := reSubmitURL.FindString(text)
	return match[10 : len(match)-1]
}

func decodeToEUCKR(s string) (string, error) {
	var buf bytes.Buffer
	wr := transform.NewWriter(&buf, korean.EUCKR.NewDecoder())
	_, err := wr.Write([]byte(s))
	if err != nil {
		return "", err
	}
	defer wr.Close()
	return buf.String(), nil
}

func findInputTags(x *html.Node, callback func(key string, value string)) {
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "input" {
			var t string
			var key string
			var val string
			for _, a := range n.Attr {
				switch a.Key {
				case "type":
					t = a.Val
				case "name":
					key = a.Val
				case "value":
					val = a.Val
				}
			}
			if t == "hidden" {
				callback(key, val)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(x)
}
