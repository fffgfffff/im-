package service

import (
	"context"
	"encoding/json"
	"fmt"
	"ginChat/models"
	"ginChat/utils"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

var (
	// 映射关系
	clientMap = make(map[int64]*models.Node, 0)
	// 读写锁
	rwLocker sync.RWMutex
	ctx      = context.Background()
	// 定义udp管道
	udpsendChan = make(chan []byte, 1024)
)

// CleanConnection 清理超时连接
func CleanConnection(param interface{}) {
	currentTime := uint64(time.Now().Unix())
	for _, node := range clientMap {
		if node.IsHeartbeatTimeOut(currentTime) {
			log.Println("心跳超时......关闭连接:", node)
			node.Conn.Close()
		}
	}
}

func ChatWithWebsocket(conn *websocket.Conn, userId int64) {
	currentTime := uint64(time.Now().Unix())
	//todo 2.获取conn
	node := &models.Node{
		Conn:          conn,
		Addr:          conn.RemoteAddr().String(), //客户端地址
		HeartbeatTime: currentTime,
		LoginTime:     currentTime,
		DataQueue:     make(chan []byte, 50),
	}
	//todo 4.发送者id和node绑定  并加锁
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()
	//todo 5.完成发送逻辑
	go sendProc(node)
	//todo 6.完成接收逻辑
	go recvProc(node)
	//todo 7.加入在线用户到缓存
	Id := strconv.FormatInt(userId, 10)
	models.SetUserOnlineInfo(Id, []byte(node.Addr),
		time.Duration(viper.GetInt("timeout.RedisOnlineTime"))*time.Hour)
	sendMsg(userId, []byte("欢迎进入聊天室"))
}

// sendProc 发送逻辑
func sendProc(node *models.Node) {
	for {
		select {
		case data := <-node.DataQueue: //并行变串行
			log.Println("[ws]sendProc data:", string(data))
			if err := node.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

// recvProc 完成接收逻辑
func recvProc(node *models.Node) {
	for {
		_, data, err := node.Conn.ReadMessage() //由onMessage发出的消息
		if err != nil {
			log.Println(err)
			return
		}
		msg := models.Message{}
		if err = json.Unmarshal(data, &msg); err != nil {
			log.Println(err)
		}
		//心跳检测 msg.Medis==-1 || msg.type=12
		if msg.Type == 12 {
			currentTime := uint64(time.Now().Unix())
			node.Heartbeat(currentTime) //重置心跳时间
		} else {

			log.Println("[ws]recvProc data:", string(data))
			//广播
			broadMsg(data)
		}
	}
}

func init() {
	//广播，服务端监听与客户端发送(udpsendChan)
	go udpSendProc()
	go udpRecvProc()
}
func sendMsg(userId int64, msg []byte) {
	rwLocker.RLock()
	node, ok := clientMap[userId] //获取节点
	rwLocker.RUnlock()

	jsonMsg := models.Message{}
	json.Unmarshal(msg, &jsonMsg) //解析msg
	jsonMsg.CreatedAt = time.Now()

	userIdStr := strconv.Itoa(int(jsonMsg.UserId))
	targetIdStr := strconv.Itoa(int(jsonMsg.TargetId))
	r, err := models.GetUserOnlineInfo(userIdStr)
	if err != nil {
		log.Println(err)
	}
	if r != "" { //???
		if ok {
			log.Println("targetID >>> userID: ", string(msg))
			node.DataQueue <- msg //输入传过来的数据
		}
	}
	var key string
	if jsonMsg.TargetId > jsonMsg.UserId { //userId=targetId
		key = models.GeneralZkey(userIdStr, targetIdStr)
	} else {
		key = models.GeneralZkey(targetIdStr, userIdStr)
	}
	fmt.Printf("%#v %v", jsonMsg, key)
	// 获取所有集合（从大到小）,不带用于排序的数字
	res, err := utils.Red.ZRevRange(ctx, key, 0, -1).Result()
	if err != nil {
		log.Println(err)
	}
	score := float64(cap(res) + 1) //每次获取都加一，保证了消息排最后
	ress, err := utils.Red.ZAdd(ctx, key, &redis.Z{score, msg}).Result()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(ress) //int64
}
func sendGroupMsg(userId, groupId int64, msg []byte) {
	userIds := models.SearchUserByGroupId(groupId)
	fmt.Println(userIds, userId)
	for i := 0; i < len(userIds); i++ {
		//排除给自己的
		if int64(userIds[i].OwnerId) != userId {
			fmt.Println(userIds[i])
			sendMsg(int64(userIds[i].OwnerId), msg)
		}
	}
}

// todo 将消息广播到局域网
func broadMsg(data []byte) {
	udpsendChan <- data
	//dispatch(data)
}

// udpSendProc 从广播通道中取出数据,完成udp数据发送
func udpSendProc() {
	//todo 使用udp协议拨号
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 125, 5),
		Port: 3000,
	})
	if err != nil {
		log.Println(err)
		return
	}
	defer con.Close()
	for {
		select {
		case data := <-udpsendChan:
			log.Println("[广播]udpSendProc data:", string(data))
			if _, err := con.Write(data); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

// udpRecvProc 广播通道,完成udp数据接收数据
func udpRecvProc() {
	//fmt.Println("start udpRecvProc")
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	if err != nil {
		log.Println("listenUdp faired:", err)
		return
	}
	defer con.Close()
	for {
		var buf [512]byte
		n, err := con.Read(buf[0:])
		if err != nil {
			log.Println("read faired:", err)
			continue
		}
		log.Println("[广播]udpRecvProc data:", string(buf[:n]))
		//后端调度逻辑处理
		dispatch(buf[0:n])
	}
}

// dispatch 后端调度逻辑处理
func dispatch(data []byte) {
	//todo 解析data为message
	msg := models.Message{}
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Println("json unmarshal faired:", err)
		return
	}
	//todo 根据type对逻辑进行处理
	switch msg.Type {
	case models.SINGLE_MSG: //私信
		fmt.Println("私信:", string(data))
		sendMsg(msg.TargetId, data) //用msg.TargetId再发一遍 前端onMessage
	case models.ROOM_MSG: //群发
		fmt.Println("群发:", string(data), msg.TargetId)
		//todo 群聊转发逻辑
		sendGroupMsg(msg.UserId, msg.TargetId, data)
	case models.HEART: //心跳
	}
}
func GetMsg(userStrA, userStrB string) ([]string, error) {
	userIdA, _ := strconv.Atoi(userStrA)
	userIdB, _ := strconv.Atoi(userStrB)
	key := ""
	if userIdA < userIdB {
		key = models.GeneralZkey(userStrA, userStrB)
	} else {
		key = models.GeneralZkey(userStrB, userStrA)
	}
	return models.GetMsgByZscore(key)
}
func CheckToken(userId int, token string) bool {
	return models.CheckToken(userId, token)
}
