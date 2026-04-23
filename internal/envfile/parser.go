package envfile

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	TypePlain = "plain"
	TypeKP    = "kp"
	TypeOP    = "op"
)

type ParsedVar struct {
	Var  string
	Type string
	Path string
}

func ParseFile(path string) ([]ParsedVar, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open env file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	out := make([]ParsedVar, 0)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, raw, ok := strings.Cut(line, "=")
		if !ok {
			return nil, fmt.Errorf("invalid env declaration at line %d: missing '='", lineNo)
		}
		key = strings.TrimSpace(key)
		if key == "" {
			return nil, fmt.Errorf("invalid env declaration at line %d: empty key", lineNo)
		}

		value := normalizeValue(strings.TrimSpace(raw))
		parsed := ParsedVar{Var: key, Type: TypePlain, Path: value}
		switch {
		case strings.HasPrefix(value, "kp://"):
			parsed.Type = TypeKP
			parsed.Path = strings.TrimPrefix(value, "kp://")
		case strings.HasPrefix(value, "op://"):
			parsed.Type = TypeOP
			parsed.Path = strings.TrimPrefix(value, "op://")
		}
		out = append(out, parsed)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read env file: %w", err)
	}
	return out, nil
}

func normalizeValue(v string) string {
	if len(v) >= 2 {
		if (v[0] == '"' && v[len(v)-1] == '"') || (v[0] == '\'' && v[len(v)-1] == '\'') {
			return v[1 : len(v)-1]
		}
	}
	return v
}
