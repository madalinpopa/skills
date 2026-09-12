package skill

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"path"

	"go.yaml.in/yaml/v3"
)

type Status string

const (
	Published Status = "published"
	Draft     Status = "draft"
)

type Metadata struct {
	Name        string
	Description string
	Status      Status
	Tags        []string
}

type Skill struct {
	Metadata
	Dir string
}

type File struct {
	Data []byte
	Mode fs.FileMode
}

const (
	fileName    = "SKILL.md"
	openAIFile  = "agents/openai.yaml"
	claudeAgent = "claude"
	claudeKey   = "x-claude"
	delimiter   = "---\n"
)

func Parse(data []byte) (Metadata, error) {
	meta, err := Describe(data)
	if err != nil {
		return Metadata{}, err
	}
	switch {
	case meta.Status == "":
		return Metadata{}, errors.New("status is required: published or draft")
	case meta.Status != Published && meta.Status != Draft:
		return Metadata{}, fmt.Errorf("status %q is not published or draft", meta.Status)
	}
	return meta, nil
}

func Describe(data []byte) (Metadata, error) {
	front, _, err := split(data)
	if err != nil {
		return Metadata{}, err
	}
	var raw struct {
		Name        string   `yaml:"name"`
		Description string   `yaml:"description"`
		Status      Status   `yaml:"status"`
		Tags        []string `yaml:"tags"`
	}
	if err := yaml.Unmarshal(front, &raw); err != nil {
		return Metadata{}, err
	}
	meta := Metadata(raw)
	switch {
	case meta.Name == "":
		return Metadata{}, errors.New("name is required")
	case meta.Description == "":
		return Metadata{}, errors.New("description is required")
	}
	return meta, nil
}

func Transform(data []byte, agent string) ([]byte, error) {
	front, body, err := split(data)
	if err != nil {
		return nil, err
	}
	var doc yaml.Node
	if err := yaml.Unmarshal(front, &doc); err != nil {
		return nil, err
	}
	if len(doc.Content) != 1 || doc.Content[0].Kind != yaml.MappingNode {
		return nil, errors.New("frontmatter must be a mapping")
	}
	mapping := doc.Content[0]

	var kept []*yaml.Node
	var claude *yaml.Node
	for i := 0; i < len(mapping.Content); i += 2 {
		key, value := mapping.Content[i], mapping.Content[i+1]
		switch key.Value {
		case "status", "tags":
			continue
		case claudeKey:
			claude = value
			continue
		}
		kept = append(kept, key, value)
	}
	if agent == claudeAgent && claude != nil {
		if claude.Kind != yaml.MappingNode {
			return nil, fmt.Errorf("%s must be a mapping", claudeKey)
		}
		kept = append(kept, claude.Content...)
	}
	mapping.Content = kept

	var out bytes.Buffer
	out.WriteString(delimiter)
	enc := yaml.NewEncoder(&out)
	enc.SetIndent(2)
	if err := enc.Encode(&doc); err != nil {
		return nil, err
	}
	if err := enc.Close(); err != nil {
		return nil, err
	}
	out.WriteString(delimiter)
	out.Write(body)
	return out.Bytes(), nil
}

func Render(fsys fs.FS, agent string) (map[string]File, error) {
	files, err := Load(fsys)
	if err != nil {
		return nil, err
	}
	return RenderFiles(files, agent)
}

func RenderFiles(files map[string]File, agent string) (map[string]File, error) {
	rendered := make(map[string]File, len(files))
	for p, file := range files {
		if p == openAIFile && agent == claudeAgent {
			continue
		}
		if p == fileName {
			data, err := Transform(file.Data, agent)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", p, err)
			}
			file.Data = data
		}
		rendered[p] = file
	}
	return rendered, nil
}

func Load(fsys fs.FS) (map[string]File, error) {
	files := map[string]File{}
	err := fs.WalkDir(fsys, ".", func(p string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		if !entry.Type().IsRegular() {
			return fmt.Errorf("%s: only regular files are supported", p)
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		data, err := fs.ReadFile(fsys, p)
		if err != nil {
			return err
		}
		files[p] = File{Data: data, Mode: info.Mode().Perm()}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func Discover(fsys fs.FS) ([]Skill, error) {
	entries, err := fs.ReadDir(fsys, "skills")
	if err != nil {
		return nil, err
	}
	var skills []Skill
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		dir := path.Join("skills", entry.Name())
		file := path.Join(dir, fileName)
		data, err := fs.ReadFile(fsys, file)
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("%s: %s is missing", dir, fileName)
		}
		if err != nil {
			return nil, err
		}
		meta, err := Parse(data)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", file, err)
		}
		if meta.Name != entry.Name() {
			return nil, fmt.Errorf("%s: name %q must match the directory name", dir, meta.Name)
		}
		skills = append(skills, Skill{Metadata: meta, Dir: dir})
	}
	return skills, nil
}

func split(data []byte) (front, body []byte, err error) {
	rest, ok := bytes.CutPrefix(data, []byte(delimiter))
	if !ok {
		return nil, nil, errors.New("missing frontmatter")
	}
	front, body, ok = bytes.Cut(rest, []byte("\n"+delimiter))
	if !ok {
		return nil, nil, errors.New("unclosed frontmatter")
	}
	return append(front, '\n'), body, nil
}
