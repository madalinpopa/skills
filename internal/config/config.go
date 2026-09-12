package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Store    Store
	Agents   map[string]Agent
	Defaults Defaults
}

type Store struct {
	Repo   string
	Branch string
}

type Agent struct {
	Project string
	Global  string
}

type Defaults struct {
	Agents []string
}

const fileName = "config.toml"

const defaultTOML = `[store]
repo = "https://github.com/madalinpopa/skills"
branch = "main"

# Which agents exist, and where each one keeps its skills.
[agents.claude]
project = ".claude/skills"
global  = "~/.claude/skills"

[agents.codex]
project = ".agents/skills"
global  = "~/.agents/skills"

[agents.gemini]
project = ".agents/skills"
global  = "~/.agents/skills"

[defaults]
agents = ["claude", "codex"]
`

func Default() Config {
	return Config{
		Store: Store{
			Repo:   "https://github.com/madalinpopa/skills",
			Branch: "main",
		},
		Agents: map[string]Agent{
			"claude": {Project: ".claude/skills", Global: "~/.claude/skills"},
			"codex":  {Project: ".agents/skills", Global: "~/.agents/skills"},
			"gemini": {Project: ".agents/skills", Global: "~/.agents/skills"},
		},
		Defaults: Defaults{Agents: []string{"claude", "codex"}},
	}
}

func Dir(xdgConfigHome, home string) string {
	if xdgConfigHome != "" {
		return filepath.Join(xdgConfigHome, "skills")
	}
	return filepath.Join(home, ".config", "skills")
}

func Load(dir string) (Config, error) {
	data, err := os.ReadFile(filepath.Join(dir, fileName))
	if errors.Is(err, os.ErrNotExist) {
		return Default(), nil
	}
	if err != nil {
		return Config{}, err
	}
	return parse(data)
}

func Init(dir string) (Config, error) {
	path := filepath.Join(dir, fileName)
	if _, err := os.Stat(path); err == nil {
		return Load(dir)
	}
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return Config{}, err
	}
	if err := writeAtomic(path, []byte(defaultTOML)); err != nil {
		return Config{}, err
	}
	return Load(dir)
}

func (a Agent) GlobalPath(home string) string {
	if rest, ok := strings.CutPrefix(a.Global, "~/"); ok {
		return filepath.Join(home, rest)
	}
	return filepath.FromSlash(a.Global)
}

func parse(data []byte) (Config, error) {
	defaults := Default()

	v := viper.New()
	v.SetConfigType("toml")
	v.SetDefault("store.repo", defaults.Store.Repo)
	v.SetDefault("store.branch", defaults.Store.Branch)
	if err := v.ReadConfig(bytes.NewReader(data)); err != nil {
		return Config{}, fmt.Errorf("parse %s: %w", fileName, err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("parse %s: %w", fileName, err)
	}
	if !v.IsSet("agents") {
		cfg.Agents = defaults.Agents
	}
	if !v.IsSet("defaults.agents") {
		cfg.Defaults.Agents = defaults.Defaults.Agents
	}
	return cfg, cfg.validate()
}

func (c Config) validate() error {
	if c.Store.Repo == "" {
		return errors.New("config: store.repo is empty")
	}
	if c.Store.Branch == "" {
		return errors.New("config: store.branch is empty")
	}
	for name, agent := range c.Agents {
		if agent.Project == "" || agent.Global == "" {
			return fmt.Errorf("config: agent %q needs both project and global paths", name)
		}
		if !insideProject(agent.Project) {
			return fmt.Errorf("config: agent %q project path %q must stay inside the project", name, agent.Project)
		}
		if !filepath.IsAbs(agent.Global) && !strings.HasPrefix(agent.Global, "~/") {
			return fmt.Errorf("config: agent %q global path %q must be absolute or start with ~/", name, agent.Global)
		}
	}
	if len(c.Defaults.Agents) == 0 {
		return errors.New("config: defaults.agents is empty")
	}
	for _, name := range c.Defaults.Agents {
		if _, ok := c.Agents[name]; !ok {
			return fmt.Errorf("config: default agent %q is not defined under [agents]", name)
		}
	}
	return nil
}

func insideProject(p string) bool {
	if filepath.IsAbs(p) || strings.HasPrefix(p, "/") {
		return false
	}
	clean := filepath.ToSlash(filepath.Clean(filepath.FromSlash(p)))
	return clean != "." && clean != ".." && !strings.HasPrefix(clean, "../")
}

func writeAtomic(path string, data []byte) error {
	tmp, err := os.CreateTemp(filepath.Dir(path), ".config-*.toml")
	if err != nil {
		return err
	}
	_, writeErr := tmp.Write(data)
	if err := errors.Join(writeErr, tmp.Close()); err != nil {
		return errors.Join(err, os.Remove(tmp.Name()))
	}
	return os.Rename(tmp.Name(), path)
}
