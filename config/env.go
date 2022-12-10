package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var Env *env

type env struct {
	BotToken        string
	RedisUrl        string
	RedisPassword   string
	RedisDomainName string
	ChannelId       string
}

func Load() *env {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		log.Fatal("Unable to load .env file")
	}

	return &env{
		BotToken:        os.Getenv("BOT_TOKEN"),
		RedisUrl:        os.Getenv("REDIS_URL"),
		ChannelId:       os.Getenv("CHANNEL_ID"),
		RedisPassword:   os.Getenv("REDIS_PASSWORD"),
		RedisDomainName: os.Getenv("REDIS_DOMAIN_NAME"),
	}
}
