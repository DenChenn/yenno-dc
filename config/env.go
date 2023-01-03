package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var Env *env

type env struct {
	BotToken             string
	RedisUrl             string
	RedisPassword        string
	RedisDomainName      string
	ChannelId            string
	AisfMasterPrivateKey string
}

func Load() *env {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		log.Fatal("Unable to load .env file")
	}

	return &env{
		BotToken:             os.Getenv("BOT_TOKEN"),
		RedisUrl:             os.Getenv("REDIS_URL"),
		ChannelId:            os.Getenv("CHANNEL_ID"),
		RedisPassword:        os.Getenv("REDIS_PASSWORD"),
		RedisDomainName:      os.Getenv("REDIS_DOMAIN_NAME"),
		AisfMasterPrivateKey: os.Getenv("AISF_MASTER_PRIVATE_KEY"),
	}
}
