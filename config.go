package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	yaml "gopkg.in/yaml.v2"
)

// Config represents a configuration file.
type Config struct {
	filename string
	format   string
	cache    *map[string]interface{}
}

var extensionsByFormat = map[string][]string{
	"json": {".json"},
	"yaml": {".yaml", ".yml"},
}

// New creates a new Config object.
func New(filename string) (*Config, error) {
	config := Config{
		filename: filename,
		format:   "",
		cache:    nil,
	}

	format, err := getConfigFormat(filepath.Ext(filename), extensionsByFormat)
	if err != nil {
		return nil, err
	}
	config.format = format

	err = config.Reload()
	if err != nil {
		return nil, err
	}

	go config.watch()

	return &config, nil
}

// Get retreives a Config option into a passed in pointer or returns an error.
func (config *Config) Get(key string, v interface{}) error {
	var val interface{}

	env, set := os.LookupEnv(key)

	if set {
		switch v.(type) {
		case *float64:
			val, err := strconv.ParseFloat(env, 64)

			if err != nil {
				return err
			}

			*v.(*float64) = val
			return nil
		case *int:
			val, err := strconv.ParseInt(env, 10, 64)

			if err != nil {
				return err
			}

			*v.(*int) = int(val)
			return nil
		default:
			val = env
		}
	} else if config.cache != nil {
		val = (*config.cache)[key]
	}

	// Cast JSON values
	switch v.(type) {
	case *string:
		if val == nil {
			val = ""
		}

		if b, ok := val.(bool); ok {
			*v.(*string) = strconv.FormatBool(b)
		} else if f, ok := val.(float64); ok {
			*v.(*string) = strconv.FormatFloat(f, 'f', -1, 64)
		} else {
			*v.(*string) = val.(string)
		}
	case *bool:
		switch val {
		case nil, 0, false, "", "0", "false":
			// falsey
			val = false
		default:
			// truthy
			val = true
		}

		*v.(*bool) = val.(bool)
	case *float64:
		if val == nil {
			val = float64(0)
		}

		if s, ok := val.(string); ok {
			pf, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return err
			}

			*v.(*float64) = pf
		} else {
			*v.(*float64) = val.(float64)
		}
	case *int:
		if val == nil {
			val = float64(0)
		}

		*v.(*int) = int(val.(float64))
	default:
		return errors.New("Type not supported")
	}

	return nil
}

// Reload clears the config cache.
func (config *Config) Reload() error {
	cache, err := primeCacheFromFile(config.filename, config.format)
	config.cache = cache

	if err != nil {
		return err
	}

	return nil
}

func (config *Config) watch() {
	l := log.New(os.Stderr, "", 0)

	// Catch SIGHUP to automatically reload cache
	sighup := make(chan os.Signal, 1)
	signal.Notify(sighup, syscall.SIGHUP)

	for {
		<-sighup
		l.Println("Caught SIGHUP, reloading config...")
		config.Reload()
	}
}

func primeCacheFromFile(file string, format string) (*map[string]interface{}, error) {
	// File exists?
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil, err
	}

	// Read file
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// Unmarshal
	var config map[string]interface{}
	switch format {
	case "json":
		if err := json.Unmarshal(raw, &config); err != nil {
			return nil, err
		}
	case "yaml":
		if err := yaml.Unmarshal(raw, &config); err != nil {
			return nil, err
		}
	}

	return &config, nil
}

func getConfigFormat(currentExt string, extensionsByFormat map[string][]string) (string, error) {
	for t, exts := range extensionsByFormat {
		for _, ext := range exts {
			if ext == currentExt {
				return t, nil
			}
		}
	}

	return "", fmt.Errorf("Unsupported config extension: '%s'", currentExt)
}
