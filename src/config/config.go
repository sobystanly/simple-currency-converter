package config

import "time"

type Config struct {
	Port        int
	LogLevel    string
	CacheExpiry time.Duration
}

var AppConfig = Config{
	Port:        8001,
	LogLevel:    "Debug",
	CacheExpiry: time.Minute * 30,
}
