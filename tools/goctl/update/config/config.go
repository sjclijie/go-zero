package config

import "github.com/sjclijie/go-zero/core/logx"

type Config struct {
	logx.LogConf
	ListenOn string
	FileDir  string
	FilePath string
}
