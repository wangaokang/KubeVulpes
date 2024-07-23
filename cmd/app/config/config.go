package config

type Config struct {
	Default DefaultOptions `yaml:"default"`
	DB      DBOptions      `yaml:"db"`
}

type DefaultOptions struct {
	Mode   string `yaml:"mode"`
	Listen int    `yaml:"listen"`
	JWTKey string `yaml:"jwt_key"`

	// 自动创建指定模型的数据库表结构，不会更新已存在的数据库表
	AutoMigrate bool `yaml:"auto_migrate"`
}

type DBOptions struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	Driver   string `yaml:"driver"`
}
