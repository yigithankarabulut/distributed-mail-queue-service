package main

import (
	"github.com/yigithankarabulut/distributed-mail-queue-service/apiserver"
	"github.com/yigithankarabulut/distributed-mail-queue-service/config"
	"log"
)

func main() {
	Conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	if err := apiserver.NewApiServer(
		apiserver.WithConfig(Conf),
		apiserver.WithServerEnv("development"),
	); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
