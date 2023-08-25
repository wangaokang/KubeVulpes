package config

type Config struct {
	//Default DefaultOptions `yaml:"default"`
	Mysql MysqlOptions `yaml:"mysql"`
}

type DefaultOptions struct {
	Listen   int    `yaml:"listen"`
	LogType  string `yaml:"log_type"`
	LogDir   string `yaml:"log_dir"`
	LogLevel string `yaml:"log_level"`
	JWTKey   string `yaml:"jwt_key"`
}

type MysqlOptions struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
}
