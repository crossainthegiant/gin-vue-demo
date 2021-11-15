package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

type User struct {
	gorm.Model
	Name string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"varchar(11);not null;unique"`
	Password string `gorm:"varchar(20);not null"`
}

func main() {
	db := InitDB()
	defer db.Close()


	r := gin.Default()
	r.POST("/api/auth/register", func(c *gin.Context) {
		//获取参数
		name := c.PostForm("name")
		telephone := c.PostForm("telephone")
		password := c.PostForm("password")

		//数据验证
		if len(telephone) != 11 {
			c.JSON(422, gin.H{
				"code": 422,
				"msg":  "手机号必须为11位",
			})
			return
		}
		if len(password) < 6 {
			c.JSON(422, gin.H{
				"code": 422,
				"msg":  "密码不能少于六位",
			})
			return
		}
		if len(name) == 0 {
			c.JSON(422, gin.H{
				"code": 422,
				"msg":  "名字不能为空",
			})
			return
		}
		log.Println(name, telephone, password)
		//判断手机号是否存在
		if isTelExist(db, telephone) {
			c.JSON(http.StatusUnprocessableEntity,gin.H{
				"code": 422,
				"msg": "手机号已存在",
			})
			return
		}
		//创建用户
		newUser := User{

			Name:     name,
			Telephone:      telephone,
			Password: password,
		}
		db.Create(&newUser)
		//返回结果
		c.JSON(200, gin.H{
			"msg": "msg received",
		})
	})
	panic(r.Run())
}

func InitDB() *gorm.DB {
	driverName := "mysql"
	host := "localhost"
	port := "3306"
	database := "ginVue"
	username := "root"
	password := "19920124q"
	charset := "utf8"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)

	db, err := gorm.Open(driverName, args)
	if err != nil {
		panic("failed to connect database,err:" + err.Error())
	}
	db.AutoMigrate(&User{})
	return db
}

func isTelExist(db *gorm.DB, telephone string) bool {
	var user User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
