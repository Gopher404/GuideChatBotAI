package config

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

type Config struct {
	DB       DBConfig       `json:"db"`
	GigaChat GigaChatConfig `json:"GigaChat"`
	TGBot    TGBotConfig    `json:"TG_bot"`
}

type DBConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DbName   string `json:"db_name"`
}

type GigaChatConfig struct {
	Token   string `json:"token"`
	Model   string `json:"model"`
	Timeout int    `json:"timeout"`
}

type TGBotConfig struct {
	Token      string `json:"token"`
	Debug      bool   `json:"debug"`
	Timeout    int    `json:"timeout"`
	MaxThreads int    `json:"max_threads"`
}

func Read(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := new(Config)

	if err := json.NewDecoder(file).Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}

var (
	typeString = reflect.TypeOf("")
	typeInt    = reflect.TypeOf(0)
	typeBool   = reflect.TypeOf(false)
)

func (cfg *Config) String() string {
	v := reflect.ValueOf(*cfg)
	return "Config:" + valueToString(v, "\t")
}

func valueToString(v reflect.Value, tab string) string {
	t := v.Type()

	switch v.Type() {
	case typeString:
		return v.String()
	case typeInt:
		return strconv.Itoa(int(v.Int()))
	case typeBool:
		return strconv.FormatBool(v.Bool())
	default:
		s := "\n"
		for i := 0; i < v.NumField(); i++ {
			s += fmt.Sprintf("%s%s: %s\n", tab, t.Field(i).Name, valueToString(v.Field(i), tab+"\t"))
		}
		return s
	}
}
