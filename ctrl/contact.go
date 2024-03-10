package ctrl

import (
	"fmt"
	"ginChat/models"
	"ginChat/service"
	"ginChat/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

var (
	contactService = service.ContactService{}
)

func SearchFriends(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	fmt.Println(userId)
	ub := contactService.SearchFriends(userId)
	//c.JSON(200, gin.H{
	//	"code":    0,
	//	"message": "查找朋友成功",
	//	"data":    ub,
	//})
	utils.RespListOk(c.Writer, ub, len(ub))
}
func CreateCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.Request.FormValue("ownerId"))
	name := c.Request.FormValue("name")
	hobby, _ := strconv.Atoi(c.Request.FormValue("hobby"))
	desc := c.Request.FormValue("desc")
	icon := c.Request.FormValue("icon")
	community := models.Community{
		Name: name, OwnerId: uint(ownerId), Type: uint(hobby),
		Desc: desc, Icon: icon,
	}
	code, msg := models.CreateCommunity(community)
	if code == 0 {
		utils.RespOk(c.Writer, nil, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}
func LoadCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.Request.FormValue("ownerId"))
	data := models.LoadCommunity(uint(ownerId))
	utils.RespListOk(c.Writer, data, len(data))
}

// AddFriend 自动添加好友
func AddFriend(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	targetId, _ := strconv.Atoi(c.Request.FormValue("targetId"))
	fmt.Println(userId, targetId)

	if err := contactService.AddFriend(userId, targetId); err != nil {
		utils.RespFail(c.Writer, err.Error())
	} else {
		utils.RespOk(c.Writer, nil, "添加成功")
	}
}
func JoinCommunity(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	comId := c.Request.FormValue("comId")
	code, msg := contactService.JoinCommunity(userId, comId)
	if code == -1 {
		utils.RespFail(c.Writer, msg)
	} else {
		utils.RespOk(c.Writer, nil, msg)
	}
}
func In(c *gin.Context) {
	c.Redirect(301, "/toChat")
}
func CreateCom(c *gin.Context) {
	c.HTML(200, "/chat/createcom.shtml", nil)
}
