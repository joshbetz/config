package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type Config struct {
	filename string
	cache    *map[string]interface{}
}

func New(filename string) *Config {
	c := Config{filename, nil}
	c.Reload()
	go c.watch()
	return &c
}

func (this *Config) Get(key string, v interface{}) error {
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
	} else if this.cache != nil {
		val = (*this.cache)[key]
	}

	// Cast JSON values
	switch v.(type) {
	case *string:
		if val == nil {
			val = ""
		}

		*v.(*string) = val.(string)
	case *bool:
		switch val {
		case nil, 0, false, "", "0", "false":
			// fasley
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

		*v.(*float64) = val.(float64)
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

func (this *Config) Reload() error {
	cache, err := primeCacheFromFile(this.filename)
	this.cache = cache

	if err != nil {
		return err
	}

	return nil
}

func (this *Config) watch() {
	// Catch SIGHUP to automatically reload cache
	sighup := make(chan os.Signal, 1)
	signal.Notify(sighup, syscall.SIGHUP)

	for {
		<-sighup
		fmt.Println("Caught SIGHUP, reloading config...")
		this.Reload()
	}
}

func primeCacheFromFile(file string) (*map[string]interface{}, error) {
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
	if err := json.Unmarshal(raw, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
