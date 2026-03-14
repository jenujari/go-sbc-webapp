package config

type Config struct {
	SweGrpcConfig SweGrpcConfig `mapstructure:"swe_grpc"`
	WebAppConfig  WebAppConfig  `mapstructure:"web_app"`
}

type SweGrpcConfig struct {
	Addr string `mapstructure:"addr"`
}

type WebAppConfig struct {
	Port    string `mapstructure:"port"`
	Appname string `mapstructure:"appname"`
}
