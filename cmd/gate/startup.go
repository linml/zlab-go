package main

import (
	"fmt"
	"math/rand"
	"time"
	"zlab/core/crypto"

	"zlab/library/nano"
	"zlab/library/nano/serialize/protobuf"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	version     = ""            // 游戏版本
	consume     = map[int]int{} // 房卡消耗配置
	forceUpdate = false
	logger      = log.WithField("component", "game")
)

//Startup ..
func Startup() {
	nano.SetLogger(log.WithField("component", "nano"))

	rand.Seed(time.Now().Unix())
	version = viper.GetString("update.version")

	heartbeat := viper.GetInt("core.heartbeat")
	if heartbeat < 5 {
		heartbeat = 5
	}
	nano.SetHeartbeatInterval(time.Duration(heartbeat) * time.Second)

	//注册所有消息
	nano.RegisterAll(SwitchMessage)

	// 加密管道
	c := crypto.NewCrypto()
	pipeline := nano.NewPipeline()
	pipeline.Inbound().PushBack(c.Inbound)
	pipeline.Outbound().PushBack(c.Outbound)

	nano.SetSerializer(protobuf.NewSerializer())

	addr := fmt.Sprintf(":%d", viper.GetInt("game-server.port"))
	nano.Listen(addr, nano.WithPipeline(pipeline))
}
