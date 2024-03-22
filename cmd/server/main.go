package main

import (
	"github.com/joho/godotenv"
	"github.com/yigithankarabulut/distributed-mail-queue-service/src/apiserver"
	"github.com/yigithankarabulut/distributed-mail-queue-service/src/config"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load .env file: %v", err)
	}
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
