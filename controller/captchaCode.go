package controller

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"nngllgjw/config"
)

func Captcha() (string, string) {
	url := config.BaseURL + config.CaptchaPath
	client := http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Unexpected status code:", resp.StatusCode)
		return "", ""
	}

	cookies := resp.Cookies()
	var jsessionid string
	for _, cookie := range cookies {
		if cookie.Name == "JSESSIONID" {
			jsessionid = cookie.Value
			break
		}
	}
	imageBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "", ""
	}

	imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)

	return jsessionid, imageBase64
}
