package models

import (
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"log"
)

const (
	SINGLE_MSG = 10
	ROOM_MSG   = 11
	HEART      = 12
)

// Message 人员关系
type Message struct {
	gorm.Model
	UserId   int64  `json:"userId,omitempty"`   //发送者
	TargetId int64  `json:"targetId,omitempty"` //接收者
	Type     int    `json:"type,omitempty"`     //发送类型:	10私聊 11群聊 12心跳
	Media    int    `json:"media,omitempty"`    //消息类型	1文字 2图片 3音频 4表情包
	Content  string `json:"content,omitempty"`  //消息内容
	Pic      string `json:"pic,omitempty"`
	Url      string `json:"url,omitempty"`
	Desc     string `json:"desc,omitempty"`
	amount   int    //其他数字统计
}

func (m *Message) TableName() string {
	return "message"
}

// Node 形成userid和node的映射关系
type Node struct {
	Conn          *websocket.Conn
	Addr          string      //客户端地址
	FirstTime     uint64      //首次连接时间
	HeartbeatTime uint64      //心跳时间
	LoginTime     uint64      //登录时间
	DataQueue     chan []byte //让conn读写，并行（两次写）变串行
}

// Heartbeat 更新用户心跳
func (node *Node) Heartbeat(currentTime uint64) {
	node.HeartbeatTime = currentTime
}

// IsHeartbeatTimeOut 用户心跳是否超时
func (node *Node) IsHeartbeatTimeOut(currentTime uint64) (timeout bool) {
	if node.HeartbeatTime+viper.GetUint64("time.HeartbeatMaxTime") <= currentTime {
		log.Println("心跳超时。。。自动下线", node)
		timeout = true
	}
	return
}
