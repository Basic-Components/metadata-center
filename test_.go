package main

import (
	"fmt"

	"github.com/Basic-Components/components_manager/connects"
	"github.com/Basic-Components/components_manager/logger"
	"github.com/Basic-Components/components_manager/models"
)

func main() {
	//url := "postgres://user:pass@host.com:5432/path?k=v"
	dburl := "postgres://postgres:postgres@localhost:5432/postgres"
	emailurl := "email://huangsizhe@xndm.tech:F7AGJLagtAUb6TRz@smtp.exmail.qq.com:465/?ssl=true&tls=false"

	logger.Init("Info", nil)
	err := connects.Email.InitFromURL(emailurl)
	if err != nil {
		fmt.Println(err)
	}
	models.Init()
	err = connects.DB.InitFromURL(dburl)
	if err != nil {
		fmt.Println(err)
	}
	defer connects.DB.Close()
	var user = models.User{}
	err = connects.DB.Cli.Model(&user).Where("name = ?", "admin").Select()
	// fmt.Println(err)
	// fmt.Println(user)
	// fmt.Println(user.Valid("admin"))
	tokenStr, err := user.SignJWT(10000)
	if err != nil {
		fmt.Println(err)
	} else {
		res, err := user.CheckJWT(tokenStr)
		if err != nil {
			fmt.Println(err)
			fmt.Println(res)
		} else {
			fmt.Println(res)
		}
	}
	// err = connects.Email.SendTest("huangsizhe@xndm.tech")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	token, err := models.UserNew("hsz", "hsz1273327@gmail.com", "qwert", "localhost:1234")
	if err != nil {
		fmt.Println(err)
	} else {
		ok, err := models.UserVerify(token)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(ok)
	}
}
