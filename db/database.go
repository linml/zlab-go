package db

import (
	"fmt"
	"zlab/core/player"
	"zlab/db/conn"
	"zlab/db/fields"
	"zlab/db/key"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

var database *Database

type Database struct {
	MongoDB *mongo.Database
	Redis   *redis.Client
}

func init() {
	database = &Database{
		MongoDB: conn.MongoDB(),
		Redis:   conn.Redis(),
	}
}

func FindPlayerByID(id string) player.Player {
	res, err := database.Redis.HMGet(key.PlayerID(id), fields.Player()...).Result()
	if err != nil {
		fmt.Println(res)
	}
	var player = new(player.Player)

}
