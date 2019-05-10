package conn

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
)

var database *mongo.Database

func init() {
	var opts = options.Client().SetHosts([]string{"localhost:27017"}).SetMaxPoolSize(100)
	var conn, err = mongo.Connect(context.Background(), opts)
	if err != nil {
		panic("connect mongo fail")
	}
	database = conn.Database("game", nil)
}

//MongoDB ..
func MongoDB() *mongo.Database {
	return database
}
