package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"go.uber.org/zap"
)

type Mysql struct {
	Host      string
	Port      int
	Database  string
	Username  string
	Password  string
	Charset   string
	ParseTime bool `toml:"parse_time"`
	Loc       string
}

type Server struct {
	IP   string
	Port string
}

type Path struct {
	FfmpegPath       string `toml:"ffmpeg_path"`
	StaticSourcePath string `toml:"static_source_path"`
}

type Config struct {
	DB     Mysql `toml:"mysql"`
	Server `toml:"server"`
	Path   `toml:"path"`
}

var Info Config

//包初始化加载时候会调用的函数
func init() {
	if _, err := toml.DecodeFile("D:\\Code\\gostudy\\demo1\\config\\config.toml", &Info); err != nil {
		panic(err)
	}
}

// DBConnectString 填充得到数据库连接字符串
func DBConnectString() string {
	arg := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%v&loc=%s",
		Info.DB.Username, Info.DB.Password, Info.DB.Host, Info.DB.Port, Info.DB.Database,
		Info.DB.Charset, Info.DB.ParseTime, Info.DB.Loc)
	zap.L().Info(arg)
	return arg
}
