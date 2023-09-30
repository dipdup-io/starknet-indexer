package config

// Config - configuration structure for indexer
type Config struct {
	Name               string    `yaml:"name" validate:"omitempty"`
	StartLevel         uint64    `yaml:"start_level" validate:"omitempty"`
	ThreadsCount       int       `yaml:"threads_count" validate:"omitempty,min=1"`
	Timeout            uint64    `yaml:"timeout" validate:"omitempty"`
	Node               *Node     `yaml:"node" validate:"omitempty"`
	Sequencer          Sequencer `yaml:"sequencer" validate:"required"`
	ClassInterfacesDir string    `yaml:"class_interfaces_dir" validate:"required,dir"`
	BridgedTokensFile  string    `yaml:"bridged_tokens_file" validate:"required,file"`
	CacheDir           string    `yaml:"cache_dir" validate:"omitempty,dir"`
	Cache              bool      `yaml:"cache" validate:"omitempty"`
}

// Node -
type Node struct {
	Url string `yaml:"url" validate:"omitempty,url"`
	Rps int    `yaml:"requests_per_second" validate:"omitempty,min=1"`
}

// Sequencer -
type Sequencer struct {
	FeederGateway string `yaml:"feeder_gateway" validate:"required,url"`
	Gateway       string `yaml:"gateway" validate:"required,url"`
	Rps           int    `yaml:"requests_per_second" validate:"omitempty,min=1"`
}
