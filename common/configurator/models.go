package configurator

type DatastorerConfig struct{}

type APIConfig struct {
	Port        uint     `mapstructure:"port"`
	TrustedCORS []string `mapstructure:"trusted_cors"`
}

type GlobalConfig struct {
	Datastorer DatastorerConfig `mapstructure:"datastorer"`
	API        APIConfig        `mapstructure:"api"`
}
