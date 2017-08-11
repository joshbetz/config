package config

import (
	"fmt"
	"testing"

	"os"
)

func BenchmarkGet(b *testing.B) {
	c, err := New("test.json")
	if err != nil {
		b.Error(err)
	}

	var v string
	for i := 0; i < b.N; i++ {
		c.Get("test", &v)
	}
}

func TestGet(t *testing.T) {
	c, err := New("test.json")
	if err != nil {
		t.Error(err)
	}

	t.Run("Get", func(t *testing.T) {
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
			t.Skip()
		})

		t.Run("should successfully cast float to string", func(t *testing.T) {
			t.Skip()
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
			t.Skip()
		})

		t.Run("should successfully cast bool to float", func(t *testing.T) {
			t.Skip()
		})

		t.Run("should successfully override with Env", func(t *testing.T) {
			var s string

			err := os.Setenv("string", "abcd")
			if err != nil {
				t.Error(err)
			}

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

		t.Run("should correctly error for nonexistant files", func(t *testing.T) {
			c2, err := New("nonexistant.json")
			if err != nil {
				fmt.Println("err nonexistant.json", err)
				return
			}
			err = c2.Reload()
			if err == nil {
				t.Error("Expected an error, got nil")
			}
		})

		t.Run("should correctly error for invalid files", func(t *testing.T) {
			c2, err := New("invalid.json")
			if err != nil {
				fmt.Println("err invalid.json", err)
				return
			}
			err = c2.Reload()
			if err == nil {
				t.Error("Expected an error, got nil")
			}
		})
	})
}
