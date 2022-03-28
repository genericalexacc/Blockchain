package main

type Config struct {
	Port int `json:"port"`
}

var globalConfig Config = Config{
	Port: 8080,
}
