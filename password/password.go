package password

import (
	"fmt"

	"github.com/m1/go-generate-password/generator"
)

type password struct {
	generator *generator.Generator
}

func New(cfg Config) (Generator, error) {
	newGenerator, err := generator.New(
		&generator.Config{
			Length:                     uint(cfg.Length),
			CharacterSet:               "",
			IncludeSymbols:             cfg.IncludeSymbols,
			IncludeNumbers:             cfg.IncludeNumbers,
			IncludeLowercaseLetters:    true,
			IncludeUppercaseLetters:    true,
			ExcludeSimilarCharacters:   true,
			ExcludeAmbiguousCharacters: true,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("generator.New err:%w", err)
	}
	return &password{
		generator: newGenerator,
	}, nil
}

func (p *password) Generate() (string, error) {
	result, err := p.generator.Generate()
	if err != nil {
		return "", fmt.Errorf("generator.Generate err: %w", err)
	}

	return *result, nil
}
