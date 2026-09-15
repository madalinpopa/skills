package tests

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"{{MODULE_PATH}}/internal/modules/auth"
	"{{MODULE_PATH}}/internal/testkit"
)

var serverURL string

func testConfig() auth.Config {
	return auth.Config{
		SigningKey:      []byte("test-signing-key-of-thirty-two-b"),
		Issuer:          "test-api",
		Audience:        "test-api",
		AccessTokenTTL:  2 * time.Minute,
		RefreshTokenTTL: 5 * time.Minute,
		CookieSecure:    false,
		CookieSameSite:  http.SameSiteLaxMode,
	}
}

func TestMain(m *testing.M) {
	if !testkit.Enabled() {
		os.Exit(m.Run())
	}

	os.Exit(run(m))
}

func run(m *testing.M) int {
	ctx := context.Background()

	env, cleanup := testkit.Start(ctx)
	defer func() {
		if err := cleanup(); err != nil {
			fmt.Println("failed to stop the test database:", err)
		}
	}()

	if err := env.Wire(ctx, auth.NewModule(testConfig(), env.Pool)); err != nil {
		panic(err)
	}

	stop := env.Serve()
	defer stop()

	serverURL = env.URL

	return m.Run()
}
