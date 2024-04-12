package main

import (
	"github.com/joho/godotenv"
	"github.com/yigithankarabulut/distributed-mail-queue-service/apiserver"
	"github.com/yigithankarabulut/distributed-mail-queue-service/config"
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
	apiserv := apiserver.New(
		apiserver.WithConfig(Conf),
		apiserver.WithServerEnv("development"),
	)
	if err := apiserv.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
