package setting

// ConfigInterface config load interface
type ConfigInterface interface {
	// Load load config
	Load() error

	// IsSet is set value
	IsSet(key string) bool

	// ReadSection read val by key,val must be a pointer
	ReadSection(key string, val interface{}) error

	// Store save config to file
	Store(path string) error
}
