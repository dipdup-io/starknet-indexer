package config

// Config - configuration structure for indexer
type Config struct {
	Name               string `yaml:"name" validate:"omitempty"`
	StartLevel         uint64 `yaml:"start_level" validate:"omitempty"`
	ThreadsCount       int    `yaml:"threads_count" validate:"omitempty,min=1"`
	Timeout            uint64 `yaml:"timeout" validate:"omitempty"`
	FeederGateway      string `yaml:"feeder_gateway" validate:"required,url"`
	Gateway            string `yaml:"gateway" validate:"required,url"`
	RequestsPerSecond  int    `yaml:"requests_per_second" validate:"omitempty,min=1"`
	ClassInterfacesDir string `yaml:"class_interfaces_dir" validate:"required,dir"`
	BridgedTokensFile  string `yaml:"bridged_tokens_file" validate:"required,file"`
}
