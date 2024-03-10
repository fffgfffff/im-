package models

import (
	"context"
	"ginChat/utils"
	"time"
)

const (
	KEY1 = "online_"
	KEY2 = "msg_"
)

func SetUserOnlineInfo(Id string, val []byte, ttl time.Duration) {
	utils.Red.Set(context.Background(), KEY1+Id, val, ttl)
}
func GetUserOnlineInfo(userIdStr string) (string, error) {
	return utils.Red.Get(context.Background(), KEY1+userIdStr).Result()
}
func GeneralZkey(userIdA, targetIdB string) string {
	return KEY2 + userIdA + "_" + targetIdB
}
func GetMsgByZscore(key string) ([]string, error) {
	return utils.Red.ZRange(context.Background(), key, 0, -1).Result()
}
