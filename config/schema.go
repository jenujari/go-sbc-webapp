package config

type Config struct {
	SweGrpcConfig SweGrpcConfig `mapstructure:"swe_grpc"`
	WebAppConfig  WebAppConfig  `mapstructure:"web_app"`
	DBConfig      DBConfig      `mapstructure:"db"`
}

type SweGrpcConfig struct {
	Addr string `mapstructure:"addr"`
}

type WebAppConfig struct {
	Port    string `mapstructure:"port"`
	Appname string `mapstructure:"appname"`
}

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"sslmode"`
}
