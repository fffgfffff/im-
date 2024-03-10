package ctrl

import (
	"fmt"
	"ginChat/models"
	"ginChat/utils"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"math/rand"
	"strconv"
)

func GetUserList(ctx *gin.Context) {
	data := models.GetUserList()
	utils.RespOk(ctx.Writer, data, "查找成功")
}
func FindUserById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Request.FormValue("id"))
	user := models.FindUserById(id)
	utils.RespOk(c.Writer, user, "查找成功")
}
func FindUserByPhoneAndPwd(c *gin.Context) {
	data := &models.UserBasic{}
	// 表单
	phone := c.Request.FormValue("mobile")
	password := c.Request.FormValue("passwd")
	user := models.FindUserByPhone(phone)

	if flag := utils.ValidPassword(password, user.Salt, user.Password); !flag {
		utils.RespFail(c.Writer, "密码不正确")
		return
	}

	pwd := utils.MakePassword(password, user.Salt)
	data = models.FindUserByPhoneAndPwd(phone, pwd)
	utils.RespOk(c.Writer, data, "查找成功")
}

func CreateUser(c *gin.Context) {
	name := c.PostForm("nickname")
	phone := c.PostForm("mobile")
	password := c.PostForm("passwd")
	repassword := c.PostForm("passwdr")
	memo := c.PostForm("memo")
	avatar := c.PostForm("avatar")

	salt := fmt.Sprintf("%06d", rand.Int31())
	if name == "" || password == "" || phone == "" {
		utils.RespFail(c.Writer, "用户名,手机号或密码不能为空")
		return
	}
	user := models.FindUserByName(name)
	if user.Name != "" {
		utils.RespFail(c.Writer, "用户名已注册")
		return
	}
	user = models.FindUserByPhone(phone)
	if user.Name != "" {
		utils.RespFail(c.Writer, "手机号已注册")
		return
	}
	if password != repassword {
		utils.RespFail(c.Writer, "两次密码不一致")
		return
	}
	user.Phone = phone
	user.Name = name
	user.Password = utils.MakePassword(password, salt)
	user.Salt = salt
	user.Memo = memo
	user.Avatar = avatar
	models.CreateUser(*user)
	utils.RespOk(c.Writer, user, "添加用户成功")
}
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.Avatar = c.PostForm("avatar")
	user.Password = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")
	if _, err := govalidator.ValidateStruct(user); err != nil {
		utils.RespFail(c.Writer, "修改参数不匹配")
		return
	}
	models.UpdateUser(user)
	utils.RespOk(c.Writer, user, "修改用户成功")
}
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	models.DelUser(user)
	utils.RespOk(c.Writer, user, "删除用户成功")
}
