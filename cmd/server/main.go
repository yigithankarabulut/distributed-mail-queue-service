package main

import (
	"github.com/yigithankarabulut/distributed-mail-queue-service/src/config"
	"log"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

}
