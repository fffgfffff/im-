package router

import (
	"ginChat/ctrl"
	"ginChat/service"
	"github.com/gin-gonic/gin"
	"os"
	"path"
)

var dir string

func Router() *gin.Engine {
	r := gin.Default()
	dir, _ = os.Getwd()

	//静态资源
	r.Static("/asset", path.Join(dir, "asset"))
	r.LoadHTMLGlob(path.Join(dir, "view/**/*"))

	//首页
	r.GET("/", service.GetIndex)
	r.GET("/index", service.GetIndex)
	r.GET("/register", service.ToRegister)
	r.GET("/toChat", service.ToChat)
	r.GET("/chat", ctrl.Chat)
	r.GET("/in", ctrl.In)
	r.GET("/createCom", ctrl.CreateCom)
	//上传文件
	r.POST("/attach/upload", ctrl.Upload)
	//用户交流模块
	contact := r.Group("/contact")
	{
		contact.POST("/searchFriends", ctrl.SearchFriends)
		contact.POST("/addFriend", ctrl.AddFriend)
		//创建群
		contact.POST("/createCommunity", ctrl.CreateCommunity)
		//群列表
		contact.POST("/loadCommunity", ctrl.LoadCommunity)
		//加群
		contact.POST("/joinCommunity", ctrl.JoinCommunity)
	}
	//用户模块
	user := r.Group("/user")
	{
		user.GET("/getUserList", ctrl.GetUserList)
		user.POST("/createUser", ctrl.CreateUser)
		user.POST("/updateUser", ctrl.UpdateUser)
		user.POST("/deleteUser", ctrl.DeleteUser)
		user.POST("/findUserByPhoneAndPwd", ctrl.FindUserByPhoneAndPwd)
		user.POST("/find", ctrl.FindUserById)
		//发送消息
		user.GET("/sendMsg", service.SendMsg)
		user.GET("/sendUserMsg", ctrl.Chat)
		user.POST("/redisMsg", ctrl.RedisMsg) //初始化消息
	}
	return r
}
