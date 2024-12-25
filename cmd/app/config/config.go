package config

type Config struct {
	DB      DBOptions      `config:"db"`
	Default DefaultOptions `config:"default"`
}

type DBOptions struct {
	Name     string `config:"name"`
	Host     string `config:"host"`
	User     string `config:"user"`
	Password string `config:"password"`
	Port     int    `config:"port"`
}

type DefaultOptions struct {
	Listen int    `config:"listen"`
	Mode   string `config:"mode"` // 此处使用 unpack 渲染不能自定义变量
	JWTKey string `config:"jwt_key"`

	// 自动创建指定模型的数据库表结构，不会更新已存在的数据库表
	AutoMigrate bool `config:"auto_migrate"`
}
