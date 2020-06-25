package main

import (
	flags "github.com/jessevdk/go-flags"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const defaultConfigFilename = "naporta-api.conf"

type config struct {
	ConfigFile    string `short:"C" long:"configfile" description:"Path to config file"`
	MongoUser     string `long:"mongo_user" description:"Mongo user"`
	MongoPassword string `long:"mongo_password" description:"Mongo password"`
	MongoServer   string `long:"mongo_server" description:"mongo server"`
	MongoDB       string `long:"mongo_db" description:"database"`
	TelegramToken string `long:"telegram_token" description:"Your telegram token"`
}

func loadConfig() (*config, error) {
	defaultConfig := config{
		ConfigFile: "naporta-api.conf",
	}

	preCfg := defaultConfig
	if _, err := flags.Parse(&preCfg); err != nil {
		return nil, err
	}

	appName := filepath.Base(os.Args[0])
	appName = strings.TrimSuffix(appName, filepath.Ext(appName))

	var configFileError error
	cfg := preCfg
	if err := flags.IniParse("naporta-api.conf", &cfg); err != nil {
		if _, ok := err.(*flags.IniError); ok {
			return nil, err
		}
		configFileError = err
	}
	if _, err := flags.Parse(&cfg); err != nil {
		return nil, err
	}

	if configFileError != nil {
		log.Printf("%v", configFileError)
	}

	return &cfg, nil
}
