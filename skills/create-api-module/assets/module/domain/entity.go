package domain

import (
	"strings"
	"time"
	"uuid"
)

type {{ENTITY_PASCAL}}UUID uuid.UUID

func New{{ENTITY_PASCAL}}UUID() {{ENTITY_PASCAL}}UUID {
	return {{ENTITY_PASCAL}}UUID(uuid.NewV7())
}

func Parse{{ENTITY_PASCAL}}UUID(s string) ({{ENTITY_PASCAL}}UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return {{ENTITY_PASCAL}}UUID{}, err
	}

	return {{ENTITY_PASCAL}}UUID(id), nil
}

func (id {{ENTITY_PASCAL}}UUID) String() string {
	return uuid.UUID(id).String()
}

type {{ENTITY_PASCAL}} struct {
	UUID      {{ENTITY_PASCAL}}UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func New{{ENTITY_PASCAL}}(name string, now time.Time) (*{{ENTITY_PASCAL}}, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrEmptyName
	}

	return &{{ENTITY_PASCAL}}{
		UUID:      New{{ENTITY_PASCAL}}UUID(),
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
