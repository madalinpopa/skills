package cmd

import (
	"errors"
	"io/fs"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"
)

func newVersionCmd(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show the CLI version and the current store commit",
		Args:  usageArgs(cobra.NoArgs),
		RunE: func(c *cobra.Command, _ []string) error {
			c.Printf("  skills  %s\n", resolveVersion(version))
			a, err := openApp(c)
			if err != nil {
				return err
			}
			_, err = os.Stat(a.store.Dir)
			if errors.Is(err, fs.ErrNotExist) {
				c.Println("  store   not initialised, run 'skills init'")
				return nil
			}
			if err != nil {
				return err
			}
			commit, err := a.store.Commit(c.Context())
			if err != nil {
				return err
			}
			c.Printf("  store   %s\n", commit)
			return nil
		},
	}
}

func resolveVersion(injected string) string {
	if injected != "" {
		return injected
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return "dev"
}
