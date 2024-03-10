package utils

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

var (
	DB  *gorm.DB
	Red *redis.Client
)

const (
	PublishKey = "websocket"
)

// Publish 发布消息到redis
func Publish(ctx context.Context, channel string, msg string) error {
	err := Red.Publish(ctx, channel, msg).Err()
	return err
}

// Subscribe 从redis获取信息
func Subscribe(ctx context.Context, channel string) (string, error) {
	sub := Red.Subscribe(ctx, channel)
	msg, err := sub.ReceiveMessage(ctx)
	return msg.Payload, err
}

func InitConfig() {
	dir, _ := os.Getwd()
	dirs := strings.Split(dir, "\\")
	viper.SetConfigFile(path.Join(strings.Join(dirs[:len(dirs)], "/"), "config/app.yml"))
	//viper.SetConfigName("app.yml")
	viper.SetConfigType("yaml")
	//viper.AddConfigPath("D:/idea_golang/gin/ginChat/config")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("config app:", viper.Get("app.yml"))
}
func InitMySQL() {
	var (
		err                        error
		name, user, password, host string
		newLoggr                   logger.Interface
	)
	//自定义日志模板,打印sql语句
	newLoggr = logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags|log.Lshortfile),
		logger.Config{
			SlowThreshold: time.Second, //慢SQL阈值
			LogLevel:      logger.Info,
			Colorful:      true,
		})
	name = viper.GetString("mysql.name")
	user = viper.GetString("mysql.user")
	password = viper.GetString("mysql.password")
	host = viper.GetString("mysql.host")
	config := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc"+
		"=Local", user, password, host, name)
	fmt.Println(config)
	if DB, err = gorm.Open(mysql.Open(config), &gorm.Config{Logger: newLoggr}); err != nil {
		fmt.Println(err)
	}
}
func InitRedis() {
	Red = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.host"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.db"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})
	pong, err := Red.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("init redis fair:", err)
	} else {
		fmt.Println("redis inited...", pong)
	}
}
