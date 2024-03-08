package controller

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/parnurzeal/gorequest"
	"net/http"
	"nngllgjw/config"
	"nngllgjw/utils"
	"strconv"
	"strings"
)

type StudentOwnScore struct {
	Coursename  string `json:"Coursename"`
	Gradepoints string `json:"Gradepoints"`
	Status      string `json:"Status"`
	Credits     string `json:"Credits"`
	Year        string `json:"Year"`
	Categories  string `json:"Categories"`
}

func GetAllStudentOwnScores(cookies *http.Cookie) ([]*StudentOwnScore, error, float64, int) {
	allScores := []*StudentOwnScore{}
	request := gorequest.New()
	resp, body, errs := request.Post(config.BaseURL+config.Personal_grades_inquiry).
		Set("Cookie", cookies.String()).
		Send("year=&term=&prop=&groupName=&para=0&sortColumn=&Submit=%E6%9F%A5%E8%AF%A2").
		End()
	if errs != nil {
		return nil, errs[0], 0, 0
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return nil, err, 0, 0
	}
	var totalCreditPoints float64
	var totalCredits float64
	var Fail int

	doc.Find("table").Each(func(i int, s *goquery.Selection) {
		if i == 1 {
			s.Find("tr").Each(func(i int, s *goquery.Selection) {
				if i != 0 {
					tds := s.Find("td")
					if utils.ExcludeGPA(strings.TrimSpace(tds.Eq(3).Text())) && strings.TrimSpace(tds.Eq(11).Text()) != "学位课" {
						credit, err := strconv.ParseFloat(strings.TrimSpace(tds.Eq(7).Text()), 64)
						if err != nil {
							fmt.Println("学分数转换错误:", err)
							return
						}
						gradePoint, err := strconv.ParseFloat(strings.TrimSpace(tds.Eq(6).Text()), 64)
						if err != nil {
							fmt.Println("绩点数转换错误:", err)
							return
						}
						if strings.TrimSpace(tds.Eq(12).Text()) == "不合格" {
							Fail++
						}
						if strings.Contains(strings.TrimSpace(tds.Eq(3).Text()), "习近平新时代中国特色社会主义思想概论") {
							tds.Eq(3).SetText("习概")
						} else if strings.Contains(strings.TrimSpace(tds.Eq(3).Text()), "马克思主义基本原理") {
							tds.Eq(3).SetText("马思")
						} else if strings.Contains(strings.TrimSpace(tds.Eq(3).Text()), "毛泽东思想和中国特色社会主义理论体系概论") {
							tds.Eq(3).SetText("毛概")
						}

						totalCreditPoints += credit * gradePoint
						totalCredits += credit
						score := &StudentOwnScore{
							Coursename:  strings.TrimSpace(tds.Eq(3).Text()),
							Gradepoints: strings.TrimSpace(tds.Eq(6).Text()),
							Credits:     strings.TrimSpace(tds.Eq(7).Text()),
							Status:      strings.TrimSpace(tds.Eq(12).Text()),
							Categories:  strings.TrimSpace(tds.Eq(11).Text()),
							Year:        strings.TrimSpace(tds.Eq(0).Text()),
						}
						allScores = append(allScores, score)
					}
				}
			})
		}
	})
	averageGPA := totalCreditPoints / totalCredits
	return allScores, nil, averageGPA, Fail
}
