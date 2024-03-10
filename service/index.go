package service

import (
	"ginChat/models"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"strconv"
)

func GetIndex(ctx *gin.Context) {
	temp, _ := template.ParseFiles("D:/idea_golang/gin/ginChat/index.html",
		"D:/idea_golang/gin/ginChat/view/chat/head.html")
	if err := temp.Execute(ctx.Writer, ""); err != nil {
		log.Println("template execute faired:", err)
	}
}
func ToChat(ctx *gin.Context) {
	ub := models.UserBasic{}
	ub.Identity = ctx.Query("token")
	id, _ := strconv.Atoi(ctx.Query("userId"))
	ub.ID = uint(id)
	ctx.HTML(200, "/chat/index.shtml", nil)
}
func ToRegister(ctx *gin.Context) {
	ctx.HTML(200, "/user/register.shtml", nil)
}
