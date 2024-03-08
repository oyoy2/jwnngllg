package router

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
	"nngllgjw/config"
	"nngllgjw/controller"
	"strconv"
	"strings"
)

func Router() {

	router := gin.Default()
	store := cookie.NewStore([]byte("secret")) // 设置存储会话的密钥
	router.Use(sessions.Sessions("session", store))
	router.LoadHTMLGlob("templates/*.html")
	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {
		jsessionid, img := controller.Captcha()
		c.HTML(200, "login.html", gin.H{
			"jsessionid": jsessionid,
			"image":      img,
			"message":    "请输入账号密码登录",
		})
	})

	router.POST("/", func(c *gin.Context) {
		// 获取表单提交的数据
		captcha := c.PostForm("captcha")
		jsessionid := c.PostForm("jsessionid")
		username := c.PostForm("username")
		password := c.PostForm("password")
		jsessionid = strings.Replace(jsessionid, " ", "", -1)
		result := controller.CheckCaptcha(jsessionid, captcha)
		if result {
			islogin := controller.Login(username, password, jsessionid, captcha)
			if islogin == "" {
				jsessionid, img := controller.Captcha()
				c.HTML(http.StatusOK, "login.html", gin.H{
					"jsessionid": jsessionid,
					"image":      img,
					"message":    "账号或密码错误！请重新输入！",
				})
			} else {
				session := sessions.Default(c)
				session.Set("loggedIn", true)
				session.Set("username", username)
				session.Set("jsessionid", islogin)
				session.Save()
				c.Redirect(http.StatusFound, "/system") // 重定向到 "/system" 页面
			}
		} else {
			jsessionid, img := controller.Captcha()
			c.HTML(http.StatusOK, "login.html", gin.H{
				"jsessionid": jsessionid,
				"image":      img,
				"message":    "验证码错误！请重新输入！",
			})
		}
	})
	router.GET("/system", func(c *gin.Context) {

		session := sessions.Default(c)
		loggedIn := session.Get("loggedIn")
		jsessionid := session.Get("jsessionid")
		if jsessionid == nil {
			c.Redirect(http.StatusFound, "/")

		}
		stringValue := jsessionid.(string)
		cookie_list := &http.Cookie{
			Name:  "JSESSIONID",
			Value: stringValue,
		}
		if loggedIn == true {

			PersonalInfo, _ := controller.GetPersonalInfo(cookie_list)
			sessionS := sessions.Default(c)
			sessionS.Save()
			c.HTML(http.StatusOK, "system.html", gin.H{
				"message":      "欢迎登录系统",
				"PersonalInfo": PersonalInfo,
			})
		} else {
			c.Redirect(http.StatusFound, "/login")
		}
	})
	router.GET("/coursearrange", func(c *gin.Context) {
		session := sessions.Default(c)
		loggedIn := session.Get("loggedIn")
		jsessionid := session.Get("jsessionid")
		stringValue := jsessionid.(string)
		cookie_list := &http.Cookie{
			Name:  "JSESSIONID",
			Value: stringValue,
		}
		if loggedIn == true {
			excel, err := controller.SaveCurrcourseAsExcel(cookie_list)
			if err != nil {
				return
			}
			c.File(excel)
		} else {
			c.Redirect(http.StatusFound, "/")
		}
	})
	router.GET("/GPA", func(c *gin.Context) {
		session := sessions.Default(c)
		loggedIn := session.Get("loggedIn")
		jsessionid := session.Get("jsessionid")
		if jsessionid == nil {
			c.Redirect(http.StatusFound, "/")

		}
		stringValue := jsessionid.(string)
		cookie_list := &http.Cookie{
			Name:  "JSESSIONID",
			Value: stringValue,
		}
		if loggedIn == true {
			sessionS := sessions.Default(c)
			GPAS, _, averageGPA, Fail := controller.GetAllStudentOwnScores(cookie_list)

			averageGPAString := strconv.FormatFloat(float64(averageGPA), 'f', 2, 64)

			sessionS.Save()
			c.HTML(http.StatusOK, "gpa.html", gin.H{
				"message":          "GPA计算",
				"GPAS":             GPAS,
				"averageGPAString": averageGPAString,
				"Fail":             Fail,
			})
		} else {
			c.Redirect(http.StatusFound, "/")
		}
	})
	router.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		c.Redirect(http.StatusFound, "/")

	})
	router.Run(":" + config.Port)
}
