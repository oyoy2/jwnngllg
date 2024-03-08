package main

import (
	"fmt"
	"nngllgjw/utils"
)

func main() {
	md5 := utils.Md5("123456")
	fmt.Println(utils.Md5(md5))
}
