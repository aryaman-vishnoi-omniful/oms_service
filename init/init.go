package appinit

import (
	"context"
	"fmt"
	"oms_service/database"
	"oms_service/redis"
	"time"

	"github.com/omniful/go_commons/config"
	goredis "github.com/omniful/go_commons/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// type contextKey string

// const DBKey contextKey = "mongoDB"

// var DB *mongo.Client

func Initialize(ctx context.Context) {
	InitializeRedis(ctx)
	InitializeDB(ctx)
	InitializeKafka(ctx)
	// return ctx
}

func InitializeDB(ctx context.Context) {
	fmt.Println("Connecting to mongo...")
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(config.GetString(ctx,"mongo.string"))

	var err error
	Db, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		return 
	}
	err = Db.Ping(ctx, nil)
	if err != nil {
		fmt.Println("Failed to ping MongoDB:", err)
		return 
	}

	fmt.Println("Successfully connected to MongoDB!")

	database.SetClient(Db)
}

func InitializeKafka(ctx context.Context) {}

func InitializeRedis(ctx context.Context) {
	redis_client:=goredis.NewClient(&goredis.Config{
		ClusterMode: config.GetBool(ctx,"redis.clusterMode"),
		Hosts:[]string{config.GetString(ctx,"redis.hosts")},
		DB:config.GetUint(ctx,"redis.db"),

	})
	fmt.Println("Initialized redis client")

	redis.SetClient(redis_client)

	

}
