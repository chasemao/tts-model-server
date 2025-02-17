package model

import (
	"context"
)

type Model interface {
	Name() string
	Fields() []*Field
	TTS(ctx context.Context, text string, args map[string]string) ([]byte, string, error)
}

type Field struct {
	Name         string    `json:"name"`
	DefaultValue string    `json:"defaultValue"`
	Options      []*Option `json:"options"`
}

type Option struct {
	Value         string   `json:"value"`
	RelatedFields []*Field `json:"relatedFields"`
}
