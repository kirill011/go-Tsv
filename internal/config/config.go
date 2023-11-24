package config

import (
	"log"
	"os"
	"time"

	"github.com/kkyr/fig"
)

type Config struct {
	Host      string        `fig: "host"`
	User      string        `fig: "user"`
	Password  string        `fig: "password"`
	Dbname    string        `fig: "dbname"`
	Port      string        `fig: "port" default:"5432"`
	Sslmode   string        `fig: "sslmode"`
	Timezone  string        `fig: "timezone" default:"Europe/Moscow"`
	DirInput  string        `fig: "dirInput"`
	DirOutput string        `fig: "dirOutput"`
	TimeSleep time.Duration `fig: "timeSleep" default: "10s"`
	ApiPort   int           `fig: "apiPort" default: "8080"`

	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func Init() *Config {
	cfg := Config{}

	cfg.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lmicroseconds)
	cfg.ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lmicroseconds)

	if err := fig.Load(&cfg,
		fig.File("config.yaml"),
		fig.Dirs("internal/config", "../internal/config")); err != nil {
		cfg.ErrorLog.Fatalln("func config.Init: ", err)
	}

	cfg.InfoLog.Println("Config loaded")
	return &cfg
}
