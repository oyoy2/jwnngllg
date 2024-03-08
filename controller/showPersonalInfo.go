package controller

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/parnurzeal/gorequest"
	"net/http"
	"nngllgjw/config"
	"strings"
)

type PersonalInfo struct {
	Username    string `json:"username"`
	RealName    string `json:"realname"`
	Department  string `json:"department"`
	Major       string `json:"major"`
	Direction   string `json:"direction"`
	StudentType string `json:"student_type"`
	level       string `json:"level"`
	Academic    string `json:"Academic"`
	Grade       string `json:"grade"`
	Class       string `json:"class"`
	IDType      string `json:"id_type"`
	IDNumber    string `json:"id_number"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	PostalCode  string `json:"postal_code"`
}

func GetPersonalInfo(cookie *http.Cookie) (*PersonalInfo, error) {
	request := gorequest.New()

	resp, body, errs := request.Get(config.BaseURL+config.Person_info).
		Set("Cookie", cookie.String()).
		End()
	if errs != nil {
		return nil, errs[0]
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	personalInfo := PersonalInfo{}

	doc.Find(".form td").Each(func(i int, s *goquery.Selection) {
		label := s.Prev().Text()
		value := s.Text()

		switch label {
		case "用户名":
			personalInfo.Username = value
		case "真实姓名":
			personalInfo.RealName = value
		case "所在院系":
			personalInfo.Department = value
		case "专业":
			personalInfo.Major = value
		case "方向":
			personalInfo.Direction = value
		case "学生类别":
			personalInfo.StudentType = value
		case "年级":
			personalInfo.Grade = value
		case "班级":
			personalInfo.Class = value
		case "证件类型":
			personalInfo.IDType = value
		case "证件号码":
			personalInfo.IDNumber = value
		case "电子邮箱":
			personalInfo.Email = value
		case "联系电话":
			personalInfo.Phone = value
		case "通讯地址":
			personalInfo.Address = value
		case "邮政编码":
			personalInfo.PostalCode = value
		}
	})
	return &personalInfo, nil
}
