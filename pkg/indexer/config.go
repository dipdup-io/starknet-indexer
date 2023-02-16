package indexer

// Config - configuration structure for indexer
type Config struct {
	Name         string `yaml:"name" validate:"omitempty"`
	StartLevel   uint64 `yaml:"start_level" validate:"omitempty"`
	ThreadsCount int    `yaml:"threads_count" validate:"omitempty,min=1"`
	Timeout      uint64 `yaml:"timeout" validate:"omitempty"`
	BaseURL      string `yaml:"base_url" validate:"required,url"`
}
