package player

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
