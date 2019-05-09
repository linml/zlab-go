package gate

import (
	"fmt"
	"math/rand"
	"time"
	"zlab/util/crypto"

	"github.com/lonng/nano"
	"github.com/lonng/nano/serialize/json"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	version = ""
	logger  = log.WithField("component", "game")
)

func Startup() {
	// set nano logger
	nano.SetLogger(log.WithField("component", "nano"))

	rand.Seed(time.Now().Unix())
	version = viper.GetString("update.version")

	heartbeat := viper.GetInt("core.heartbeat")
	if heartbeat < 5 {
		heartbeat = 5
	}
	nano.SetHeartbeatInterval(time.Duration(heartbeat) * time.Second)

	// 房卡消耗配置
	// csm := viper.GetString("core.consume")
	// SetCardConsume(csm)

	logger.Infof("当前游戏服务器版本: %s, 当前心跳时间间隔: %d秒", version, heartbeat)
	logger.Info("game service starup")

	// register game handler
	nano.Register(defaultManager)
	nano.Register(defaultRoomManager)

	// 加密管道
	c := crypto.NewCrypto()
	pipeline := nano.NewPipeline()
	pipeline.Inbound().PushBack(c.Inbound)
	pipeline.Outbound().PushBack(c.Outbound)

	nano.SetSerializer(json.NewSerializer())
	addr := fmt.Sprintf(":%d", viper.GetInt("game-server.port"))
	nano.Listen(addr, nano.WithPipeline(pipeline))
}
