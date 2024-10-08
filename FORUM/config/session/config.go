package session

import (
	"FORUM/config/config"
	"strings"
)
type driver string

const (
	Memory driver = "memory"
	Redis  driver = "redis"
)

var defaultName = "FORUM-session"

type sessionConfig struct {
	Driver string
	Name   string
}

func getConfig() sessionConfig {

	wc := sessionConfig{}
	wc.Driver = string(Memory)
	if config.Config.IsSet("session.driver") {
		wc.Driver = strings.ToLower(config.Config.GetString("session.driver"))
	}

	wc.Name = defaultName
	if config.Config.IsSet("session.name") {
		wc.Name = strings.ToLower(config.Config.GetString("session.name"))
	}

	return wc
}

func getRedisConfig() redisConfig {
	Info := redisConfig{
		Host:     "127.0.0.1",
		Port:     "6379",
		DB:       0,
		Password: "",
	}
	if config.Config.IsSet("redis.host") {
		Info.Host = config.Config.GetString("redis.host")
	}
	if config.Config.IsSet("redis.port") {
		Info.Port = config.Config.GetString("redis.port")
	}
	if config.Config.IsSet("redis.db") {
		Info.DB = config.Config.GetInt("redis.db")
	}
	if config.Config.IsSet("redis.pass") {
		Info.Password = config.Config.GetString("redis.pass")
	}
	return Info
}