package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tkanos/gonfig"
	"sceyt_task/pkg/logging"
	"sync"
)

const (
	ServerAddr = ""
	ServerPort = "8080"

	RedisHost    = "redis_db"
	RedisPort    = "6379"
	RedisDb      = 0
	RedisExpires = 300

	LogConfigFileName = "logConfig"
	ServerConfigPath  = "./properties"
	DbConfigPath      = "./properties/dbConfig.yml"
)

const (
	GroupPath   = "/"
	DeletePath  = "delete/"
	AddPath     = "add/"
	UpdatePath  = "update/"
	SearchPath  = "search"
	SwaggerPath = "/swagger/*any"
)

// Configuration wraps all the configs variables required by the auth service
type Configuration struct {
	Username     string
	Password     string
	Address      string
	ProtoVersion int
	Keyspace     string
	CQLVersion   string
}

var instance *logging.Configuration
var logOnce sync.Once

// GetLogConfiguration reads log configuration from a config file
func GetLogConfiguration() *logging.Configuration {
	logOnce.Do(func() {
		viper.SetConfigName(LogConfigFileName)
		viper.AddConfigPath(ServerConfigPath)
		err := viper.ReadInConfig() // Find and read the uploaderConfig file
		if err != nil {             // Handle errors reading the uploaderConfig file
			logrus.Error(err)
		}
		err = viper.Unmarshal(&instance)
		if err != nil {
			logrus.Error(err)
		}
	})
	return instance
}

var dbConfig *Configuration
var dbOnce sync.Once

// LoadConfig get sql connection parameters
func LoadConfig() *Configuration {
	dbOnce.Do(func() {
		config := &Configuration{}
		err := gonfig.GetConf(DbConfigPath, config)
		if err != nil {
			logrus.Error("An error was generated while reading the database config file.")
			return
		}
		dbConfig = config
	})
	return dbConfig
}
