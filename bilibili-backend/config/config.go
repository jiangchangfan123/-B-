package config

//引入 viper 库，Go 生态最流行的配置管理工具，支持 YAML/JSON/TOML/环境变量等多种格式。
import (
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	DB       DBConfig       `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	MinIO    MinIOConfig    `mapstructure:"minio"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
	Redis    RedisConfig    `mapstructure:"redis"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DBConfig struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Charset  string `mapstructure:"charset"`
}

type JWTConfig struct {
	Secret    string `mapstructure:"secret"`
	ExpiresIn int    `mapstructure:"expires_in"`
}

type MinIOConfig struct {
	Endpoint      string `mapstructure:"endpoint"`
	AccessKey     string `mapstructure:"access_key"`
	SecretKey     string `mapstructure:"secret_key"`
	UseSSL        bool   `mapstructure:"use_ssl"`
	BucketVideos  string `mapstructure:"bucket_videos"`
	BucketCovers  string `mapstructure:"bucket_covers"`
	BucketAvatars string `mapstructure:"bucket_avatars"`
}

type RabbitMQConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Queue    string `mapstructure:"queue"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// 全局配置变量
var C *Config

func Init() error {
	viper.SetConfigName("config")   // 配置文件名：config（不带后缀）
	viper.SetConfigType("yaml")     // 文件格式：YAML
	viper.AddConfigPath("./config") //先找 ./config/config.yaml
	viper.AddConfigPath(".")        // 再找 ./config.yaml（兜底）

	// 默认值
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("jwt.secret", "nebula-tv-secret-key-2026")
	viper.SetDefault("jwt.expires_in", 7200)

	//读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	//创建空配置对象
	C = &Config{}
	// 把 viper 里的数据映射到 C 的各个字段
	return viper.Unmarshal(C)
}
