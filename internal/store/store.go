package store

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing/fstest"
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
	if err := s.checkClean(ctx); err != nil {
		return Result{}, err
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

func (s Store) Preview(ctx context.Context) (Result, error) {
	if err := s.checkClean(ctx); err != nil {
		return Result{}, err
	}
	local, err := s.Commit(ctx)
	if err != nil {
		return Result{}, err
	}
	out, err := s.git(ctx, "ls-remote", "--quiet", "origin", "refs/heads/"+s.Branch)
	if err != nil {
		return Result{}, err
	}
	remote, _, ok := strings.Cut(out, "\t")
	if !ok {
		return Result{}, fmt.Errorf("branch %s does not exist on %s", s.Branch, s.Repo)
	}
	return Result{Old: local, New: remote}, nil
}

func (s Store) Commit(ctx context.Context) (string, error) {
	return s.git(ctx, "rev-parse", "HEAD")
}

func (s Store) Tree(ctx context.Context, commit string) (fs.FS, error) {
	listing, err := output(ctx, s.Dir, nil, "ls-tree", "-r", "-z", commit)
	if err != nil {
		return nil, fmt.Errorf("commit %s: %w", commit, err)
	}
	tree := fstest.MapFS{}
	var oids []string
	for entry := range bytes.SplitSeq(bytes.TrimSuffix(listing, []byte{0}), []byte{0}) {
		meta, name, ok := bytes.Cut(entry, []byte{'\t'})
		if !ok {
			continue
		}
		fields := strings.Fields(string(meta))
		if len(fields) != 3 || fields[1] != "blob" {
			continue
		}
		tree[string(name)] = &fstest.MapFile{Mode: blobMode(fields[0])}
		oids = append(oids, fields[2]+" "+string(name))
	}
	if len(oids) == 0 {
		return tree, nil
	}
	stdin := strings.NewReader(strings.Join(oids, "\n") + "\n")
	blobs, err := output(ctx, s.Dir, stdin, "cat-file", "--batch=%(objectname) %(objectsize) %(rest)")
	if err != nil {
		return nil, err
	}
	for len(blobs) > 0 {
		header, rest, ok := bytes.Cut(blobs, []byte{'\n'})
		if !ok {
			return nil, fmt.Errorf("git cat-file: unexpected output %q", header)
		}
		fields := strings.SplitN(string(header), " ", 3)
		if len(fields) != 3 {
			return nil, fmt.Errorf("git cat-file: unexpected output %q", header)
		}
		size, err := strconv.Atoi(fields[1])
		if err != nil || len(rest) < size+1 {
			return nil, fmt.Errorf("git cat-file: unexpected output %q", header)
		}
		tree[fields[2]].Data = bytes.Clone(rest[:size])
		blobs = rest[size+1:]
	}
	return tree, nil
}

func blobMode(mode string) fs.FileMode {
	switch mode {
	case "100755":
		return 0o755
	case "120000":
		return fs.ModeSymlink | 0o777
	default:
		return 0o644
	}
}

func (s Store) clone(ctx context.Context) error {
	if err := os.MkdirAll(filepath.Dir(s.Dir), 0o750); err != nil {
		return err
	}
	_, err := run(ctx, "", "clone", "--quiet", "--branch", s.Branch, "--single-branch", s.Repo, s.Dir)
	return err
}

func (s Store) checkClean(ctx context.Context) error {
	if err := s.check(ctx); err != nil {
		return err
	}
	status, err := s.git(ctx, "--no-optional-locks", "status", "--porcelain")
	if err != nil {
		return err
	}
	if status != "" {
		return fmt.Errorf("store %s has local changes; inspect them with 'git -C %s status'", s.Dir, s.Dir)
	}
	return nil
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
	out, err := output(ctx, dir, nil, args...)
	return strings.TrimSpace(string(out)), err
}

func output(ctx context.Context, dir string, stdin io.Reader, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir
	cmd.Stdin = stdin
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
