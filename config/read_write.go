package config

import (
	"github.com/pkg/errors"
	"github.com/tryffel/market/modules/Error"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// ReadConfig read config file
func ReadConfig(location string) (*Config, error) {

	c := &Config{}
	data, err := ioutil.ReadFile(location)
	if err != nil {
		return c, &Error.Error{Code: Error.Einvalid,
			Err: errors.Wrap(err, "failed to read configuration file")}
	}

	err = yaml.Unmarshal(data, c)
	if err != nil {
		return c, &Error.Error{Code: Error.Einvalid,
			Err: errors.Wrap(err, "failed to parse configuration file")}
	}

	err = c.SaveFile(location)
	return c, nil
}

//SaveFile saves configuration file
func (c *Config) SaveFile(location string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return &Error.Error{Code: Error.Einternal,
			Err: errors.Wrap(err, "failed to parse configuration file")}
	}

	err = ioutil.WriteFile(location, data, 0600)
	if err != nil {
		return &Error.Error{Code: Error.Einvalid,
			Err: errors.Wrap(err, "failed to save configuration file")}
	}
	return nil
}
