package tests

import (
	"context"
	"fmt"
	"os"
	"testing"

	"{{MODULE_PATH}}/internal/modules/{{MODULE_NAME}}"
	"{{MODULE_PATH}}/internal/testkit"
)

var serverURL string

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

	if err := env.Wire(ctx, {{MODULE_NAME}}.NewModule({{MODULE_NAME}}.Config{}, env.Pool)); err != nil {
		panic(err)
	}

	stop := env.Serve()
	defer stop()

	serverURL = env.URL

	return m.Run()
}
