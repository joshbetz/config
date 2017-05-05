# config

[![Build Status](https://travis-ci.org/joshbetz/config.svg?branch=master)](https://travis-ci.org/joshbetz/config) [![](https://godoc.org/github.com/joshbetz/config?status.svg)](http://godoc.org/github.com/joshbetz/config)


A small configuration library for Go that parses environment variables, JSON
files, and reloads automatically on `SIGHUP`.

## Example

```go
func main() {
	c := config.New("config.json")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var value string
		c.Get("value", &value)
		fmt.Fprintf(w, "Value: %s", value)
	})

	http.ListenAndServe(":3000", nil)
}
```

![Reload config on SIGHUP](http://i.imgur.com/6H8b6zy.gif)

## API

```go
func New(file string) *Config
```

Constructor that initializes a Config object and sets up the SIGHUP watcher.

```go
func (config *Config) Get(key string, v interface{}) error
```

Takes the path to a JSON file, the name of the configuration option, and a
pointer to the variable where the config value will be stored. `v` can be a
pointer to a string, bool, or float64.

```go
func (config *Config) Reload()
```

Reloads the config. Happens automatically on `SIGHUP`.
