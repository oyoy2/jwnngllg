package controller

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/parnurzeal/gorequest"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"net/http"
	"nngllgjw/config"
	"strings"
)

func ListLeft(cookie *http.Cookie) map[string]string {
	request := gorequest.New()

	resp, body, errs := request.Get(config.BaseURL+config.List).
		Set("Cookie", cookie.String()).
		End()
	if errs != nil {
		log.Fatal(errs)
	}
	defer resp.Body.Close()

	utf8Body, err := DecodeGBKToUTF8([]byte(body))
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(utf8Body))
	if err != nil {
		log.Fatal(err)
	}

	links := make(map[string]string)
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if href != "#" {
			if exists {
				spanContent := s.Find("span").Text()
				links[href] = spanContent
			}
		}
	})

	return links
}

func DecodeGBKToUTF8(data []byte) (string, error) {
	reader := transform.NewReader(strings.NewReader(string(data)), simplifiedchinese.GBK.NewDecoder())
	utf8Data, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(utf8Data), nil
}
