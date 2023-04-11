package setting

// Option config option
type Options struct {
	configFile string
}

// Option for ConfigOption
type Option func(*Options)

// WithConfigFile set config filename
func WithConfigFile(configFile string) Option {
	return func(c *Options) {
		c.configFile = configFile
	}
}
