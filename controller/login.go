package controller

import (
	"github.com/parnurzeal/gorequest"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/url"
	"nngllgjw/config"
	"strings"
)

func Login(username string, password string, cookie string, CaptchaCode string) string {
	request := gorequest.New()

	params := url.Values{}
	params.Set("j_username", username)
	params.Set("j_password", password)
	params.Set("j_captcha", CaptchaCode)
	// 将查询参数添加到URL中
	url := config.BaseURL + config.LoginPath + "?" + params.Encode()
	resp, body, _ := request.Get(url).
		Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7").
		Set("Accept-Encoding", "gzip, deflate").
		Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6").
		Set("Cookie", "JSESSIONID="+cookie).
		Set("Host", "jw.glutnn.cn").
		Set("Proxy-Connection", "keep-alive").
		Set("Referer", "http://jw.glutnn.cn/academic/common/security/affairLogin.jsp").
		Set("Upgrade-Insecure-Requests", "1").
		Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36 Edg/119.0.0.0").
		End()

	utf8Body, _ := ioutil.ReadAll(transform.NewReader(strings.NewReader(body), simplifiedchinese.GBK.NewDecoder()))
	body = string(utf8Body)
	if !strings.Contains(body, "综合教务管理系统") {
		return ""
	}
	Cookies := resp.Request.Cookies()
	var jsessionid string
	for Cookie := range Cookies {
		if Cookies[Cookie].Name == "JSESSIONID" {
			jsessionid = Cookies[Cookie].Value
		}
	}
	return jsessionid
}
