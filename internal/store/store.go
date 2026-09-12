package store

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/madalinpopa/skills/internal/skill"
)

type Store struct {
	Dir    string
	Repo   string
	Branch string
}

type Result struct {
	Old string
	New string
}

func (s Store) Init(ctx context.Context) error {
	_, err := os.Stat(s.Dir)
	if errors.Is(err, os.ErrNotExist) {
		return s.clone(ctx)
	}
	if err != nil {
		return err
	}
	return s.check(ctx)
}

func (s Store) Sync(ctx context.Context) (Result, error) {
	if err := s.check(ctx); err != nil {
		return Result{}, err
	}
	status, err := s.git(ctx, "status", "--porcelain")
	if err != nil {
		return Result{}, err
	}
	if status != "" {
		return Result{}, fmt.Errorf("store %s has local changes; inspect them with 'git -C %s status'", s.Dir, s.Dir)
	}
	old, err := s.Commit(ctx)
	if err != nil {
		return Result{}, err
	}
	if _, err = s.git(ctx, "fetch", "--quiet", "origin", s.Branch); err != nil {
		return Result{}, err
	}
	if _, err = s.git(ctx, "merge", "--ff-only", "--quiet", "FETCH_HEAD"); err != nil {
		return Result{}, fmt.Errorf("store %s cannot fast-forward to origin/%s; inspect it with 'git -C %s status': %w", s.Dir, s.Branch, s.Dir, err)
	}
	head, err := s.Commit(ctx)
	if err != nil {
		return Result{}, err
	}
	return Result{Old: old, New: head}, nil
}

func (s Store) Commit(ctx context.Context) (string, error) {
	return s.git(ctx, "rev-parse", "HEAD")
}

func (s Store) Files(ctx context.Context, commit, dir string) (map[string]skill.File, error) {
	archive, err := output(ctx, s.Dir, "archive", "--format=tar", commit, dir)
	if err != nil {
		return nil, err
	}
	files := map[string]skill.File{}
	reader := tar.NewReader(bytes.NewReader(archive))
	for {
		header, err := reader.Next()
		if errors.Is(err, io.EOF) {
			return files, nil
		}
		if err != nil {
			return nil, err
		}
		rel, ok := strings.CutPrefix(header.Name, dir+"/")
		if !ok || header.Typeflag != tar.TypeReg {
			continue
		}
		data, err := io.ReadAll(reader)
		if err != nil {
			return nil, err
		}
		files[rel] = skill.File{Data: data, Mode: header.FileInfo().Mode().Perm()}
	}
}

func (s Store) clone(ctx context.Context) error {
	if err := os.MkdirAll(filepath.Dir(s.Dir), 0o750); err != nil {
		return err
	}
	_, err := run(ctx, "", "clone", "--quiet", "--branch", s.Branch, "--single-branch", s.Repo, s.Dir)
	return err
}

func (s Store) check(ctx context.Context) error {
	if _, err := os.Stat(filepath.Join(s.Dir, ".git")); err != nil {
		return fmt.Errorf("%s exists but is not a Git clone of %s", s.Dir, s.Repo)
	}
	origin, err := s.git(ctx, "remote", "get-url", "origin")
	if err != nil {
		return err
	}
	if origin != s.Repo {
		return fmt.Errorf("%s is a clone of %s, not of the configured store %s", s.Dir, origin, s.Repo)
	}
	branch, err := s.git(ctx, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return err
	}
	if branch != s.Branch {
		return fmt.Errorf("%s is on branch %s, not the configured branch %s", s.Dir, branch, s.Branch)
	}
	return nil
}

func (s Store) git(ctx context.Context, args ...string) (string, error) {
	return run(ctx, s.Dir, args...)
}

func run(ctx context.Context, dir string, args ...string) (string, error) {
	out, err := output(ctx, dir, args...)
	return strings.TrimSpace(string(out)), err
}

func output(ctx context.Context, dir string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		message := strings.TrimSpace(stderr.String())
		if message == "" {
			message = err.Error()
		}
		return nil, fmt.Errorf("git %s: %s", args[0], message)
	}
	return stdout.Bytes(), nil
}
