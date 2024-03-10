package ctrl

import (
	"ginChat/service"
	"ginChat/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

// Chat :发送者id 接收者id 消息类型 发送的内容 发送类型
func Chat(c *gin.Context) {
	r := c.Request
	w := c.Writer
	//todo 1.获取参数并验证 token 合法性
	query := r.URL.Query()
	token := query.Get("token")
	Id := query.Get("userId")
	userId, _ := strconv.ParseInt(Id, 10, 64)
	isValida := service.CheckToken(int(userId), token)
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			//token校验
			return isValida
		},
	}).Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	service.ChatWithWebsocket(conn, userId)
}

func RedisMsg(c *gin.Context) {
	userStrA := c.Request.FormValue("userIdA")
	userStrB := c.Request.FormValue("userIdB")
	msgs, err := service.GetMsg(userStrA, userStrB)
	if err != nil {
		utils.RespListFail(c.Writer, nil, nil)
	} else {
		utils.RespListOk(c.Writer, msgs, len(msgs))
	}
}
