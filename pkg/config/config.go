package config

import (
	"flag"
	"github.com/AntonBorzenko/RestrictedPassportService/utils"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port int `yml:"port" env:"PORT" env-default:"80"`
	SslKey string `yml:"ssl_key" env:"SSL_KEY" env-default:"" env-description:"Path to SSL Key file"`
	SslCert string `yml:"ssl_cert" env:"SSL_CHAIN" env-default:"" env-description:"Path to SSL Cert file"`

	DBFileUrl       string `yml:"db_url" env:"DB_URL" env-default:"http://guvm.mvd.ru/upload/expired-passports/list_of_expired_passports.csv.bz2"`
	DBFile          string `yml:"db_file" env:"DB_FILE" env-default:"db.sqlite"`
	DBBatchSize     int    `yml:"db_batch_size" env:"DB_BATCH_SIZE" env-default:"2048" env-description:"Count of rows in one DB insert"`
	DBUpdateVerbose bool   `yml:"db_update_verbose" env:"DB_UPDATE_VERBOSE" env-default:"true"  env-description:"Verbose output on DB Update"`
}

var Cfg = Config{}

func Load() {
	err := cleanenv.ReadEnv(&Cfg)
	if err != nil {
		panic(err)
	}
}

func LoadFromConfig(filename string) {
	err := cleanenv.ReadConfig(filename, &Cfg)
	if err == nil {
		return
	}
	if err.Error()  != "config file parsing error: EOF" {
		panic(err)
	}
	Load()
}

func getConfigName(args []string) string {
	fset := flag.NewFlagSet("config", flag.ContinueOnError)
	configPath := fset.String("cfg", "", "path to config file")
	_ = fset.Parse(args)

	if *configPath != "" {
		return *configPath
	}
	if utils.FileExists("config.yml") {
		return "config.yml"
	}
	return ""
}

func LoadConfigFromArgs(args []string) {
	if configName := getConfigName(args); configName != "" {
		LoadFromConfig(configName)
	} else {
		Load()
	}
}
