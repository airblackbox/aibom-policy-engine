// Copyright 2024 Nostalgic Skin Co.
// SPDX-License-Identifier: AGPL-3.0-or-later

package engine

import (
	"encoding/json"
	"fmt"
	"os"
)

// RuleBuilder is a function that creates a Rule from params.
type RuleBuilder func(sev Severity, params map[string]any) (Rule, error)

// Registry maps rule type names to builders.
type Registry struct {
	builders map[string]RuleBuilder
}

// NewRegistry creates an empty rule registry.
func NewRegistry() *Registry {
	return &Registry{builders: make(map[string]RuleBuilder)}
}

// Register adds a rule builder.
func (r *Registry) Register(typ string, builder RuleBuilder) {
	r.builders[typ] = builder
}

// LoadPolicyFile loads a policy from a JSON file using the registry.
func (r *Registry) LoadPolicyFile(path string) (*Policy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read policy: %w", err)
	}
	return r.ParsePolicyJSON(data)
}

// ParsePolicyJSON parses a policy from JSON.
func (r *Registry) ParsePolicyJSON(data []byte) (*Policy, error) {
	var pf PolicyFile
	if err := json.Unmarshal(data, &pf); err != nil {
		return nil, fmt.Errorf("parse policy: %w", err)
	}

	p := &Policy{Name: pf.Name}
	for _, rd := range pf.Rules {
		sev := SeverityHigh
		if rd.Severity != "" {
			sev = Severity(rd.Severity)
		}

		builder, ok := r.builders[rd.Type]
		if !ok {
			return nil, fmt.Errorf("unknown rule type: %s", rd.Type)
		}
		rule, err := builder(sev, rd.Params)
		if err != nil {
			return nil, fmt.Errorf("build rule %s: %w", rd.Type, err)
		}
		p.Rules = append(p.Rules, rule)
	}
	return p, nil
}
