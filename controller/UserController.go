package controller

import (
	"github.com/crossainthegiant/gin-vue/common"
	"github.com/crossainthegiant/gin-vue/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

func Register(c *gin.Context) {
	DB := common.GetDB()
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
	if isTelExist(DB, telephone) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "手机号已存在",
		})
		return
	}
	//创建用户
	newUser := model.User{

		Name:      name,
		Telephone: telephone,
		Password:  password,
	}
	DB.Create(&newUser)
	//返回结果
	c.JSON(200, gin.H{
		"msg": "msg received",
	})
}

func isTelExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
