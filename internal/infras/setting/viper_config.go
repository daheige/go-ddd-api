package setting

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// New create a config interface.
func New(opts ...Option) ConfigInterface {
	c := &viperConfig{
		vp: viper.New(),
	}

	conf := &Options{
		configFile: "./app.yaml",
	}

	for _, o := range opts {
		o(conf)
	}

	c.configFile = conf.configFile

	return c
}

type viperConfig struct {
	vp         *viper.Viper
	configFile string
}

// Load load config
func (c *viperConfig) Load() error {
	if c.configFile == "" {
		return errors.New("config file is empty")
	}

	configDir, err := filepath.Abs(filepath.Dir(c.configFile))
	if err != nil {
		return err
	}

	// file ext
	ext := strings.TrimPrefix(filepath.Ext(c.configFile), ".")
	if ext == "" {
		ext = "yaml" // default yaml file
	}

	file := filepath.Base(c.configFile)
	c.vp.SetConfigName(file)
	c.vp.AddConfigPath(configDir)
	c.vp.SetConfigType(ext)
	err = c.vp.ReadInConfig()
	if err != nil {
		return err
	}

	return nil
}

// IsSet is set value
func (c *viperConfig) IsSet(key string) bool {
	return c.vp.IsSet(key)
}

// ReadSection read val by key,val must be a pointer
func (c *viperConfig) ReadSection(key string, val interface{}) error {
	return c.vp.UnmarshalKey(key, val)
}

// Store save config to file
func (c *viperConfig) Store(path string) error {
	return c.vp.WriteConfigAs(path)
}
