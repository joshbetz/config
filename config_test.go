package config

import (
	"os"
	"testing"
)

func BenchmarkGet(b *testing.B) {
	c, _ := New("test.json")

	var v string
	for i := 0; i < b.N; i++ {
		c.Get("test", &v)
	}
}

func TestGet(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		t.Run("should successfully create new config", func(t *testing.T) {
			c, err := New("test.json")
			if err != nil {
				t.Error(err)
			} else if c.format != "json" {
				t.Errorf("Expected format to be 'json', got '%v'", c.format)
			}

			c2, err := New("test.yaml")
			if err != nil {
				t.Error(err)
			} else if c2.format != "yaml" {
				t.Errorf("Expected format to be 'yaml', got '%v'", c2.format)
			}
		})

		t.Run("should return unsupported type error", func(t *testing.T) {
			_, err := New("foobar.exe")

			if err == nil || err.Error() != "Unsupported config extension: '.exe'" {
				t.Error("Expected an unsupported config extension error")
			}
		})

		t.Run("should correctly error for nonexistant files", func(t *testing.T) {
			_, err := New("nonexistant.json")
			if err == nil {
				t.Error("Expected an error, got nil")
			}
		})

		t.Run("should correctly error for invalid files", func(t *testing.T) {
			_, err := New("invalid.json")
			if err == nil {
				t.Error("Expected an error, got nil")
			}
		})
	})

	t.Run("Get", func(t *testing.T) {
		var testFiles = []string{"test.json", "test.yaml"}

		for _, testFile := range testFiles {
			c, _ := New(testFile)
			t.Run("should successfully retrieve string", func(t *testing.T) {
				var s string

				err := c.Get("string", &s)
				if err != nil {
					t.Error(err)
				}

				if s != "asdf" {
					t.Errorf("Expected 'asdf', got '%v'", s)
				}
			})

			t.Run("should successfully cast empty value to string", func(t *testing.T) {
				var s string

				err := c.Get("nonexistant", &s)
				if err != nil {
					t.Error(err)
				}

				if s != "" {
					t.Errorf("Expected '', got '%v'", s)
				}
			})

			t.Run("should successfully cast bool to string", func(t *testing.T) {
				var s string

				err := c.Get("bool", &s)
				if err != nil {
					t.Error(err)
				}

				if s != "true" {
					t.Errorf("Expect 'true', got '%v'", s)
				}
			})

			t.Run("should successfully cast float to string", func(t *testing.T) {
				var s string

				err := c.Get("float64", &s)
				if err != nil {
					t.Error(err)
				}

				if s != "64.4" {
					t.Errorf("Expect '64.4', got '%v'", s)
				}
			})

			t.Run("should successfully retrieve bool", func(t *testing.T) {
				var b bool

				err := c.Get("bool", &b)
				if err != nil {
					t.Error(err)
				}

				if b != true {
					t.Errorf("Expected 'true', got '%v'", b)
				}
			})

			t.Run("should successfully cast string to bool", func(t *testing.T) {
				var b bool

				err := c.Get("string", &b)
				if err != nil {
					t.Error(err)
				}

				if b != true {
					t.Errorf("Expected 'true', got '%v'", b)
				}
			})

			t.Run("should successfully cast empty value to bool", func(t *testing.T) {
				var b bool

				err := c.Get("nonexistant", &b)
				if err != nil {
					t.Error(err)
				}

				if b != false {
					t.Errorf("Expected 'false', got '%v'", b)
				}
			})

			t.Run("should successfully cast float to bool", func(t *testing.T) {
				var b bool

				// truthy
				err := c.Get("float64", &b)
				if err != nil {
					t.Error(err)
				}

				if b != true {
					t.Errorf("Expected 'true', got '%v'", b)
				}

				// falsey
				err = c.Get("falsey", &b)
				if err != nil {
					t.Error(err)
				}

				if b != false {
					t.Errorf("Expected 'false', got '%v'", b)
				}
			})

			t.Run("should successfully retrieve float64", func(t *testing.T) {
				var f float64

				err := c.Get("float64", &f)
				if err != nil {
					t.Error(err)
				}

				if f != 64.4 {
					t.Errorf("Expected '64.4', got '%v'", f)
				}
			})

			t.Run("should successfully cast empty value to float64", func(t *testing.T) {
				var f float64

				err := c.Get("nonexistant", &f)
				if err != nil {
					t.Error(err)
				}

				if f != 0 {
					t.Errorf("Expected '0', got '%v'", f)
				}
			})

			t.Run("should successfully cast string to float", func(t *testing.T) {
				var f float64

				err := c.Get("string2", &f)
				if err != nil {
					t.Error(err)
				}

				if f != 13.37 {
					t.Errorf("Expect '13.37', got '%v'", f)
				}

				err = c.Get("string3", &f)
				if err != nil {
					t.Error(err)
				}

				if f != 0.0 {
					t.Errorf("Expect '0.0', got '%v'", f)
				}
			})

			t.Run("should successfully override with Env", func(t *testing.T) {
				var s string

				err := os.Setenv("string", "abcd")
				if err != nil {
					t.Error(err)
				}
				defer os.Unsetenv("string")

				err = c.Get("string", &s)
				if err != nil {
					t.Error(err)
				}

				if s != "abcd" {
					t.Errorf("Expected 'abcd', got '%v'", s)
				}
			})

			t.Run("should correctly cast env to float64", func(t *testing.T) {
				var f float64

				err := os.Setenv("float64-2", "50")
				if err != nil {
					t.Error(err)
				}
				defer os.Unsetenv("float64-2")

				err = c.Get("float64-2", &f)
				if err != nil {
					t.Error(err)
				}

				if f != 50 {
					t.Errorf("Expected '50', got '%v'", f)
				}
			})

			t.Run("should correctly cast env to int", func(t *testing.T) {
				var i int

				err := os.Setenv("int", "50")
				if err != nil {
					t.Error(err)
				}
				defer os.Unsetenv("int")

				err = c.Get("int", &i)
				if err != nil {
					t.Error(err)
				}

				if i != 50 {
					t.Errorf("Expected '50', got '%v'", i)
				}
			})

			t.Run("should successfully cast float to int", func(t *testing.T) {
				var i int

				err := c.Get("float64", &i)
				if err != nil {
					t.Error(err)
				}

				if i != 64 {
					t.Errorf("Expected '64', got '%v'", i)
				}
			})

			t.Run("should correctly cast empty value to int", func(t *testing.T) {
				var i int

				err := c.Get("nonexistant", &i)
				if err != nil {
					t.Error(err)
				}

				if i != 0 {
					t.Errorf("Expected '0', got '%v'", i)
				}
			})

			t.Run("should correctly error for invalid types", func(t *testing.T) {
				var r rune

				err := c.Get("float64", &r)
				if err == nil {
					t.Error("Expected an error, got nil")
				}
			})
		}
	})
}
