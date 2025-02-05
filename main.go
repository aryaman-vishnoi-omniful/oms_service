package main

import (
	// "net/http"
	"context"
	"fmt"
	"oms_service/init"
	"oms_service/router"
	"time"

	// "github.com/gin-gonic/gin"
	// "github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/shutdown"
	// "github.com/omniful/go_commons/worker/configs"
)

func main() {
	err := config.Init(time.Second * 10)
	if err != nil {
		log.Panicf("Error while initialising config, err: %v", err)
		panic(err)
	}
	ctx, err := config.TODOContext()
	if err != nil {
		log.Panicf("Error while getting context from config, err: %v", err)
		panic(err)
	}
	// SetupServer(ctx)
	// appinit.Initialize(ctx)
	// ctx :=context.TODO()

	// appinit.InitializeDB(ctx)
	appinit.Initialize(ctx)
	SetupServer(context.TODO())

}

func SetupServer(ctx context.Context) {
	server := http.InitializeServer(config.GetString(ctx, "server.port"), 10*time.Second, 10*time.Second, 70*time.Second)
	err := router.Initialize(ctx, server)
	if err != nil {
		log.Errorf(err.Error())
		panic(err)
	}
	fmt.Println("helllo")
	log.Infof("Starting server on port" + config.GetString(ctx, "server.port"))
	err=server.StartServer(config.GetString(ctx, "service.name"))
	if err!=nil{
		log.Errorf(err.Error())
		panic(err)

	}
	fmt.Println("helllo")
	<-shutdown.GetWaitChannel()

}
