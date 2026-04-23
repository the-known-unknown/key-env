package secrets

import (
	"fmt"
	"strings"

	"github.com/asimmittal/key-env/internal/envfile"
	"github.com/asimmittal/key-env/internal/vault"
)

type LoadedVar struct {
	Var   string
	Type  string
	Path  string
	Value string
}

type Loader struct {
	providers map[string]vault.Provider
}

func NewLoader(providers map[string]vault.Provider) *Loader {
	return &Loader{providers: providers}
}

func (l *Loader) Load(parsed []envfile.ParsedVar) ([]LoadedVar, error) {
	out := make([]LoadedVar, 0, len(parsed))
	for _, item := range parsed {
		if item.Type == envfile.TypePlain {
			out = append(out, LoadedVar{
				Var:   item.Var,
				Type:  item.Type,
				Path:  item.Path,
				Value: item.Path,
			})
			continue
		}

		provider, ok := l.providers[item.Type]
		if !ok {
			return nil, fmt.Errorf("no provider configured for type %q (var %s)", item.Type, item.Var)
		}

		value, err := provider.Resolve(item.Path)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve %s (%s://%s): %w", item.Var, item.Type, item.Path, err)
		}
		out = append(out, LoadedVar{
			Var:   item.Var,
			Type:  item.Type,
			Path:  item.Path,
			Value: value,
		})
	}
	return out, nil
}

func MergeWithCurrentEnv(loaded []LoadedVar, current []string) []string {
	merged := make(map[string]string, len(current)+len(loaded))
	for _, item := range current {
		k, v, ok := strings.Cut(item, "=")
		if !ok {
			continue
		}
		merged[k] = v
	}
	for _, item := range loaded {
		merged[item.Var] = item.Value
	}
	out := make([]string, 0, len(merged))
	for k, v := range merged {
		out = append(out, k+"="+v)
	}
	return out
}
