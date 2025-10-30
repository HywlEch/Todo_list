package config

import (
	"github.com/spf13/viper"
)

// Config 结构体用于映射配置文件中的所有配置项
type Config struct {
	Database DBConfig
	Server   ServerConfig
	JWT      JWTConfig
	Redis    RedisConfig
}

// DBConfig 结构体用于映射 database 部分的配置
type DBConfig struct {
	User     string
	Password string
	DBName   string
	Host     string
	Port     int
	SSLMode  string
}

// ServerConfig 结构体用于映射 server 部分的配置
type ServerConfig struct {
	Port string
}

//JWTConfig 结构体用于映射 jwt 部分的配置
type JWTConfig struct {
	Secret 			string
    ExpiresInHours	int
}

//RedisConfig 结构体用于映射 redis 部分的配置
type RedisConfig struct {
	Addr 		string
	Password 	string
	DB 			int
}
// LoadConfig 从 config.yaml 文件加载配置
func LoadConfig() (config Config, err error) {
	// 设置配置文件的名称和类型
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	// 添加配置文件的搜索路径（. 表示当前项目根目录）
	viper.AddConfigPath(".")

	// 读取配置文件
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	// 将读取到的配置信息反序列化到 Config 结构体中
	err = viper.Unmarshal(&config)
	return
}