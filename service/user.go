package service

import (
	"fmt"
	"ginChat/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

// 防止跨域站点伪造请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {
	conn, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	MsgHandler(conn, c)
}

func MsgHandler(conn *websocket.Conn, c *gin.Context) {
	msg, err := utils.Subscribe(c, utils.PublishKey) //获取信息
	if err != nil {
		log.Println(err)
		return
	}

	if err := utils.Publish(c, utils.PublishKey, "nihao"); err != nil {
		log.Println(err)
		return
	}

	tm := time.Now().Format("2006-01-02 15:04:05")
	m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
	if err = conn.WriteMessage(1, []byte(m)); err != nil {
		log.Println(err)
	}
}
