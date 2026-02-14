// Copyright 2024 Nostalgic Skin Co.
// SPDX-License-Identifier: AGPL-3.0-or-later

package rules

import (
	"fmt"

	"github.com/nostalgicskinco/aibom-policy-engine/pkg/engine"
)

// DefaultRegistry returns a registry with all built-in rules registered.
func DefaultRegistry() *engine.Registry {
	r := engine.NewRegistry()

	r.Register("deny-model", func(sev engine.Severity, params map[string]any) (engine.Rule, error) {
		name, _ := params["model_name"].(string)
		reason, _ := params["reason"].(string)
		if name == "" {
			return nil, fmt.Errorf("model_name required")
		}
		return &DenyModel{ModelName: name, Reason: reason, Sev: sev}, nil
	})

	r.Register("require-model-version", func(sev engine.Severity, params map[string]any) (engine.Rule, error) {
		name, _ := params["model_name"].(string)
		if name == "" {
			return nil, fmt.Errorf("model_name required")
		}
		return &RequireModelVersion{ModelName: name, Sev: sev}, nil
	})

	r.Register("allowed-providers", func(sev engine.Severity, params map[string]any) (engine.Rule, error) {
		var providers []string
		if ps, ok := params["providers"].([]any); ok {
			for _, p := range ps {
				if s, ok := p.(string); ok {
					providers = append(providers, s)
				}
			}
		}
		return &AllowedProviders{Providers: providers, Sev: sev}, nil
	})

	r.Register("deny-tool", func(sev engine.Severity, params map[string]any) (engine.Rule, error) {
		name, _ := params["tool_name"].(string)
		reason, _ := params["reason"].(string)
		if name == "" {
			return nil, fmt.Errorf("tool_name required")
		}
		return &DenyTool{ToolName: name, Reason: reason, Sev: sev}, nil
	})

	r.Register("max-models", func(sev engine.Severity, params map[string]any) (engine.Rule, error) {
		maxF, _ := params["max"].(float64)
		return &MaxModels{Max: int(maxF), Sev: sev}, nil
	})

	r.Register("max-tools", func(sev engine.Severity, params map[string]any) (engine.Rule, error) {
		maxF, _ := params["max"].(float64)
		return &MaxTools{Max: int(maxF), Sev: sev}, nil
	})

	r.Register("require-all-models-versioned", func(sev engine.Severity, _ map[string]any) (engine.Rule, error) {
		return &RequireAllModelsVersioned{Sev: sev}, nil
	})

	r.Register("deny-external-endpoints", func(sev engine.Severity, params map[string]any) (engine.Rule, error) {
		var hosts []string
		if hs, ok := params["allowed_hosts"].([]any); ok {
			for _, h := range hs {
				if s, ok := h.(string); ok {
					hosts = append(hosts, s)
				}
			}
		}
		return &DenyExternalEndpoints{AllowedHosts: hosts, Sev: sev}, nil
	})

	return r
}
