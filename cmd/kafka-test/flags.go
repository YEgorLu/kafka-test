package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/YEgorLu/kafka-test/internal/config"
	"github.com/YEgorLu/kafka-test/internal/logger"
	"github.com/YEgorLu/kafka-test/internal/util"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type paramBinding[T any] struct {
	from string
	to   T
}

type fileParamBinding struct {
	paramBinding[*string]
	ext string
}

var filePathParams = []fileParamBinding{
	{paramBinding: paramBinding[*string]{"LOGS_PATH", &config.App.LogsPath}, ext: ".json"},
}

var folderPathParams = []paramBinding[*string]{
	{"MIGRATIONS_FOLDER", &config.DB.MigrationsFolderPath},
}

var stringsParams = []paramBinding[*string]{
	{"PORT", &config.App.Port},
	{"HOST_ADDR", &config.App.ServerHostAddr},
	{"DB_PROVIDER", &config.DB.Provider},
	{"DB_URL", &config.DB.Url},
	{"DB_USER", &config.DB.User},
	{"DB_PASSWORD", &config.DB.Password},
	{"DB_DBNAME", &config.DB.DbName},
	{"KAFKA_MESSAGES_TOPIC", &config.Mq.MessagesTopic},
	{"KAFKA_ERR_MESSAGES_TOPIC", &config.Mq.ErrMessagesTopic},
}

var boolsParams = []paramBinding[*bool]{
	{"DB_APPEND_MIGRATIONS", &config.DB.AppendMigrations},
	{"DB_FORCE_RECREATE", &config.DB.ForceRecreate},
}

var stringArrayParams = []paramBinding[*[]string]{
	{"KAFKA_BROKERS", &config.Mq.Brokers},
}

func init() {
	parseFlags()
	parseEnvironment()
	parseFileConfig()
}

func parseFlags() {
	flag.Func("config-path", "Path to .env file", checkFile(&config.App.EnvPath, ".env"))
	flag.Func("log-path", "Path to .json file for logs", checkFile(&config.App.LogsPath, ".json"))
	flag.Func("migrations-path", "Absolute or relative to project root path to migratinos folder", checkFolder(&config.DB.MigrationsFolderPath, true))

	flag.StringVar(&config.App.Port, "port", "8080", "Port to run server")
	flag.StringVar(&config.App.ServerHostAddr, "host-addr", "localhost", "Host address on server where app deployed")
	flag.StringVar(&config.DB.Provider, "db-provider", "postgres", "Database provider name")
	flag.StringVar(&config.DB.DbName, "db-name", "TimeTracker", "Database provider name")
	flag.StringVar(&config.DB.Url, "db-url", "localhost", "Url to database server")
	flag.StringVar(&config.DB.Port, "db-port", "5432", "Url to database server")
	flag.StringVar(&config.DB.User, "db-user", "admin", "Username of database user")
	flag.StringVar(&config.DB.Password, "db-password", "admin", "Password of database user")
	flag.BoolVar(&config.DB.AppendMigrations, "db-append-migrations", false, "Delete database and recreate from migrations")

	flag.StringVar(&config.Mq.MessagesTopic, "msg-topic", config.Mq.MessagesTopic, "Topic name for messages")
	flag.StringVar(&config.Mq.ErrMessagesTopic, "err-msg-topic", config.Mq.ErrMessagesTopic, "Topic name for messages, that failed to process")

	flag.Func("brokers", "kafka brokers urls, example: localhost:9092,localhost:9093", func(s string) error {
		config.Mq.Brokers = parseStringsArray(s)
		return nil
	})

	flag.Func("log-level", "log level for logger: 'panic', 'fatal', 'error', 'warn' | 'warning', 'info', 'debug', 'trace'", func(s string) error {
		logLevel, err := logrus.ParseLevel(s)
		if err != nil {
			return err
		}
		fmt.Println(logLevel)
		config.App.LogLevel = logLevel
		return nil
	})

	flag.Parse()
}

func parseEnvironment() {
	for _, v := range filePathParams {
		if envValue := os.Getenv(v.from); envValue != "" {
			checkFile(v.to, v.ext)(envValue)
		}
	}

	for _, v := range folderPathParams {
		if envValue := os.Getenv(v.from); envValue != "" {
			checkFolder(v.to, true)(envValue)
		}
	}

	for _, v := range stringsParams {
		if envValue := os.Getenv(v.from); envValue != "" {
			*v.to = envValue
		}
	}

	for _, v := range boolsParams {
		if envValue, err := strconv.ParseBool(os.Getenv(v.from)); err == nil {
			*v.to = envValue
		}
	}

	for _, v := range stringArrayParams {
		*v.to = parseStringsArray(os.Getenv(v.from))
	}

	if logLevel, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL")); err == nil {
		config.App.LogLevel = logLevel
	}

	fmt.Println(config.App, config.DB, config.Mq)
}

func parseFileConfig() {
	if config.App.EnvPath == "" {
		return
	}
	curDir, err := os.Getwd()
	if err != nil {
		logger.Get().Error(err)
		return
	}
	envMap, err := godotenv.Read(path.Join(curDir, config.App.EnvPath))
	if err != nil {
		logger.Get().Error(errors.Join(errors.New("invalid .env path provided"), err))
		return
	}

	for _, v := range filePathParams {
		if envValue, ok := envMap[v.from]; ok {
			checkFile(v.to, v.ext)(envValue)
		}
	}

	for _, v := range folderPathParams {
		if envValue, ok := envMap[v.from]; ok {
			checkFolder(v.to, true)(envValue)
		}
	}

	for _, v := range stringsParams {
		if envValue, ok := envMap[v.from]; ok {
			*v.to = envValue
		}
	}

	for _, v := range boolsParams {
		if envValue, ok := envMap[v.from]; ok {
			if envBoolValue, err := strconv.ParseBool(envValue); err != nil {
				logger.Get().Error(err)
			} else {
				*v.to = envBoolValue
			}
		}
	}

	if envValue, ok := envMap["LOG_LEVEL"]; ok {
		if logLevel, err := logrus.ParseLevel(envValue); err == nil {
			config.App.LogLevel = logLevel
		}
	}
}

func checkFile(toPastePath *string, ext string) func(string) error {
	return func(s string) error {
		exists, err := util.FileExists(s)
		if err != nil {
			panic(err)
		}
		if !exists {
			return os.ErrNotExist
		}
		if !util.FileHasExt(s, ext) {
			return newExtensionError(s, ext)
		}
		writable, err := util.FileWritable(s)
		if err != nil {
			panic(err)
		}
		if !writable {
			return newNotWritableError(s)
		}
		*toPastePath = s
		return nil
	}
}

func checkFolder(toPastePath *string, isAbsolute bool) func(string) error {
	return func(s string) error {
		if path.IsAbs(s) != isAbsolute {
			return newPathError(s, isAbsolute)
		}

		v, err := os.Stat(s)
		if err != nil {
			return err
		}
		if !v.IsDir() {
			return newNotFolderError(s)
		}

		*toPastePath = s
		return nil
	}
}

func newExtensionError(path string, wantExt string) error {
	return fmt.Errorf("invalid file extension, want: %s, path: %s", wantExt, path)
}

func newNotWritableError(path string) error {
	return fmt.Errorf("file not writable, path: %s", path)
}

func newPathError(s string, isAbsolute bool) error {
	absoluteTxt := "absolute"
	if !isAbsolute {
		absoluteTxt = "relative"
	}
	return fmt.Errorf("path %s is not %s", s, absoluteTxt)
}

func newNotFolderError(s string) error {
	return fmt.Errorf("%s is not a folder", s)
}

func parseStringsArray(s string) []string {
	itemsRaw := strings.Trim(s, " ,.")
	items := strings.Split(itemsRaw, ",")
	return items
}
