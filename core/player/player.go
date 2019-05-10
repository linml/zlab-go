package player

import (
	"encoding/json"
	"zlab/core/message"
	"zlab/db/fields"

	"go.mongodb.org/mongo-driver/bson"
)

type Player struct {
	ID        string
	NickName  string
	Avatar    string
	OpenID    string
	UnionID   string
	GameLevel uint8 //游戏角色类型 机器人 代理 普通玩家

	Money int64 //金币数量

	RoomID string //当前所在房间
	State  uint8  //玩家状态
}

func (p *Player) RedisLoad(values []interface{}) *Player {
	if fs := fields.Player(); len(values) != len(fs) {
		return p
	}
	p.ID = values[0].(string)

	return p
}

func (p *Player) AddMoney(money int64) error {

	var update = bson.M{"$inc": bson.M{
		"money": money,
	}}
	var cmd, _ = json.Marshal(update)
	message.Producer.DbAsync(cmd)
	return nil
}
