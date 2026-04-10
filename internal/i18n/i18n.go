package i18n

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Store holds translations for all locales
type Store struct {
	translations map[string]map[string]any
}

// New creates a Store from raw JSON bytes for each locale.
// enJSON is the base (fallback), ptBRJSON is layered on top for pt-br.
func New(enJSON, ptBRJSON []byte) (*Store, error) {
	var en map[string]any
	if err := json.Unmarshal(enJSON, &en); err != nil {
		return nil, fmt.Errorf("parse en.json: %w", err)
	}

	var ptBR map[string]any
	if err := json.Unmarshal(ptBRJSON, &ptBR); err != nil {
		return nil, fmt.Errorf("parse pt-br.json: %w", err)
	}

	merged := deepMerge(en, ptBR)

	return &Store{
		translations: map[string]map[string]any{
			"en":    en,
			"pt-br": merged,
		},
	}, nil
}

// Get retrieves a translation value by locale and dot-separated key path.
// Returns an empty string if the key is not found.
func (s *Store) Get(locale, path string) string {
	tree, ok := s.translations[locale]
	if !ok {
		tree = s.translations["en"]
	}
	return dotGet(tree, path)
}

// TFunc returns a translation function bound to a specific locale.
func (s *Store) TFunc(locale string) func(string) string {
	return func(path string) string {
		return s.Get(locale, path)
	}
}

// dotGet traverses a nested map using dot-path notation.
func dotGet(m map[string]any, path string) string {
	parts := strings.SplitN(path, ".", 2)
	val, ok := m[parts[0]]
	if !ok {
		return ""
	}
	if len(parts) == 1 {
		if s, ok := val.(string); ok {
			return s
		}
		return fmt.Sprintf("%v", val)
	}
	sub, ok := val.(map[string]any)
	if !ok {
		return ""
	}
	return dotGet(sub, parts[1])
}

// deepMerge merges src into a copy of base (src values win on conflicts).
func deepMerge(base, src map[string]any) map[string]any {
	result := make(map[string]any, len(base))
	for k, v := range base {
		result[k] = v
	}
	for k, sv := range src {
		if bv, ok := result[k]; ok {
			bMap, bIsMap := bv.(map[string]any)
			sMap, sIsMap := sv.(map[string]any)
			if bIsMap && sIsMap {
				result[k] = deepMerge(bMap, sMap)
				continue
			}
		}
		result[k] = sv
	}
	return result
}
