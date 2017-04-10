# config

[![Build Status](https://travis-ci.org/joshbetz/config.svg?branch=master)](https://travis-ci.org/joshbetz/config)

A small configuration library for Go that parses environment variables, JSON
files, and reloads automatically on `SIGHUP`.

## Example

```go
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var value string
		config.Get("config.json", "value", &value)
		fmt.Fprintf(w, "Value: %s", value)
	})

	http.ListenAndServe(":3000", nil)
}
```

![Reload config on SIGHUP](http://i.imgur.com/6H8b6zy.gif)

## API

```go
func Get(file, key string, v interface{}) error
```

Takes the path to a JSON file, the name of the configuration option, and a
pointer to the variable where the config value will be stored. `v` can be a
pointer to a string, bool, or float64.

```go
func Reload()
```

Reloads the config. Happens automatically on `SIGHUP`.
