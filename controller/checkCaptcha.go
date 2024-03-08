package controller

import (
	"io/ioutil"
	"net/http"
	"nngllgjw/config"
)

func CheckCaptcha(cookie string, captcha string) bool {
	CheckPath := config.BaseURL + config.CheckPath + captcha
	client := http.Client{}
	req, err := http.NewRequest("GET", CheckPath, nil)
	if err != nil {
		// 处理错误
		return false
	}
	req.AddCookie(&http.Cookie{Name: "JSESSIONID", Value: cookie})
	resp, err := client.Do(req)
	if err != nil {
		// 处理错误
		return false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// 处理错误
		return false
	}

	result := string(body)
	if result == "true" {
		return true
	} else {
		return false
	}
}
