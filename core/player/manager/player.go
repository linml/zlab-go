package manager

import (
	"github.com/go-redis/redis"
)

type PlayerManager struct {
	cache redis.Client
}

func (m *PlayerManager) Player(id string) {
	m.cache.HGet("key", "field")
}

func (m *PlayerManager) SetPlayer() {

}
