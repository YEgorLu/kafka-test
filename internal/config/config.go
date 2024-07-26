package config

import "github.com/sirupsen/logrus"

var App = struct {
	Port           string
	ServerHostAddr string
	EnvPath        string
	LogsPath       string
	LogLevel       logrus.Level
}{
	Port:           "8080",
	ServerHostAddr: "localhost",
	EnvPath:        "../.env",
	LogsPath:       "../logs.json",
	LogLevel:       logrus.ErrorLevel,
}

var DB = struct {
	Provider             string
	Url                  string
	Port                 string
	User                 string
	Password             string
	DbName               string
	AppendMigrations     bool
	MigrationsFolderPath string
	ForceRecreate        bool
}{
	Provider:             "postgres",
	Url:                  "localhost",
	Port:                 "5432",
	User:                 "admin",
	Password:             "admin",
	DbName:               "dbname",
	AppendMigrations:     false,
	MigrationsFolderPath: "./migrations",
	ForceRecreate:        false,
}

var Mq = struct {
	Brokers          []string
	MessagesTopic    string
	ErrMessagesTopic string
}{
	Brokers:          []string{"localhost:9092"},
	MessagesTopic:    "messages",
	ErrMessagesTopic: "err-messages",
}
