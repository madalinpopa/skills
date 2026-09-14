package testkit

import (
	"os"
	"testing"
)

func Enabled() bool {
	return os.Getenv("INTEGRATION") == "true"
}

func RequireIntegration(t *testing.T) {
	t.Helper()

	if !Enabled() {
		t.Skip("set INTEGRATION=true to run integration tests")
	}
}
