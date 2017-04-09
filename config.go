package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
)

var cache map[string]map[string]interface{}

func Get(file, key string, v interface{}) error {
	var val interface{}

	val, set := os.LookupEnv(key)

	if !set {
		if cache[file] == nil {
			if err := primeCacheFromFile(file); err != nil {
				return err
			}
		}

		val = cache[file][key]
	}

	switch v.(type) {
	case *string:
		if val == nil {
			val = ""
		}

		*v.(*string) = val.(string)
	case *bool:
		switch val {
		case nil, 0, false, "", "0", "false":
			val = false
		default:
			val = true
		}

		*v.(*bool) = val.(bool)
	case *float64:
		if val == nil {
			val = float64(0)
		}

		*v.(*float64) = val.(float64)
	default:
		return errors.New("Type not supported")
	}

	return nil
}

func Reload() {
	// Empty the cache
	cache = make(map[string]map[string]interface{})
}

func primeCacheFromFile(file string) error {
	// File exists?
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return err
	}

	// Read file
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	// Unmarshal
	var config map[string]interface{}
	if err := json.Unmarshal(raw, &config); err != nil {
		return err
	}

	// prime cache
	cache[file] = config
	return nil
}

func init() {
	Reload()

	// Catch SIGHUP to automatically reload cache
	sighup := make(chan os.Signal, 1)
	signal.Notify(sighup, syscall.SIGHUP)

	go func() {
		for sig := range sighup {
			switch sig {
			case syscall.SIGHUP:
				fmt.Println("Caught SIGHUP, reloading config...")
				Reload()
			}
		}
	}()
}
