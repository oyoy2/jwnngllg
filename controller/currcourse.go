package controller

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/parnurzeal/gorequest"
	"github.com/tealeg/xlsx"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"math/rand"
	"net/http"
	"nngllgjw/config"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func SaveCurrcourseAsExcel(cookie *http.Cookie) (string, error) {
	file := xlsx.NewFile()

	request := gorequest.New()
	resp, body, _ := request.Get(config.BaseURL+config.Current_semester_schedule).
		Set("Cookie", cookie.String()).
		End()
	defer resp.Body.Close()
	utf8Body, _ := ioutil.ReadAll(transform.NewReader(strings.NewReader(body), simplifiedchinese.GBK.NewDecoder()))
	body = string(utf8Body)
	pattern := `window\.open\('(\.\./\.\./manager/coursearrange/showTimetable\.do\?id=\d+&yearid=\d+&termid=\d+&timetableType=STUDENT&sectionType=BASE)'\)`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(body)
	url := matches[1][6:]
	resp, body, _ = request.Get(config.BaseURL+url).
		Set("Cookie", cookie.String()).
		End()
	defer resp.Body.Close()
	utf8Body, _ = ioutil.ReadAll(transform.NewReader(strings.NewReader(body), simplifiedchinese.GBK.NewDecoder()))
	body = string(utf8Body)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return "", err
	}
	table := doc.Find("#timetable")
	tableData := [][]string{}

	table.Find("tr").Each(func(rowIdx int, row *goquery.Selection) {
		rowData := []string{}

		row.Find("th, td").Each(func(colIdx int, cell *goquery.Selection) {
			rowData = append(rowData, cell.Text())
		})

		tableData = append(tableData, rowData)
	})
	sheet, err := file.AddSheet("课表")
	if err != nil {
		return "", err
	}
	for _, rowData := range tableData {
		row := sheet.AddRow()
		for _, cellData := range rowData {
			cell := row.AddCell()
			cell.Value = cellData
		}
	}
	rand.Seed(time.Now().UnixNano())
	fileName := "课表_" + strconv.Itoa(rand.Intn(10000)) + ".xlsx"
	err = file.Save(fileName)
	if err != nil {
		return "", err
	}

	return fileName, nil
}
