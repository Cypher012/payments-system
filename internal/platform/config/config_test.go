package config_test

import (
	"strings"
	"testing"
	"time"

	"github.com/Cypher012/payments-system/internal/platform/config"
)

func TestEnvReportsEveryProblem(t *testing.T) {
	vars := map[string]string{
		"PORT":    "eighty",
		"MODE":    "fast",
		"NAME":    " ",
		"TIMEOUT": "not-a-duration",
	}

	env := config.New(func(key string) (string, bool) {
		v, ok := vars[key]
		return v, ok
	})

	port := env.Int("PORT", 8080)
	env.OneOf("MODE", "safe", "safe", "strict")
	env.Required("NAME")
	timeout := env.Duration("TIMEOUT", 5*time.Second)

	if port != 8080 || timeout != 5*time.Second {
		t.Errorf("got port %d and timeout %s, want the defaults", port, timeout)
	}

	err := env.Err()
	if err == nil {
		t.Fatal("want an error, got nil")
	}
	for _, key := range []string{"PORT", "MODE", "NAME", "TIMEOUT"} {
		if !strings.Contains(err.Error(), key+":") {
			t.Errorf("error doesn't mention %s:\n%v", key, err)
		}
	}
}

func TestEnvValidConfiguration(t *testing.T) {
	vars := map[string]string{
		"PORT":    "9000",
		"MODE":    "strict",
		"NAME":    "wallet",
		"TIMEOUT": "10s",
	}

	env := config.New(func(key string) (string, bool) {
		v, ok := vars[key]
		return v, ok
	})

	port := env.Int("PORT", 8080)
	mode := env.OneOf("MODE", "safe", "safe", "strict")
	name := env.Required("NAME")
	timeout := env.Duration("TIMEOUT", 5*time.Second)

	if port != 9000 {
		t.Errorf("port = %d, want 9000", port)
	}

	if mode != "strict" {
		t.Errorf("mode = %q, want strict", mode)
	}

	if name != "wallet" {
		t.Errorf("name = %q, want wallet", name)
	}

	if timeout != 10*time.Second {
		t.Errorf("timeout = %s, want 10s", timeout)
	}

	if err := env.Err(); err != nil {
		t.Errorf("unexpected configuration error: %v", err)
	}
}
