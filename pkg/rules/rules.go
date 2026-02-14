// Copyright 2024 Nostalgic Skin Co.
// SPDX-License-Identifier: AGPL-3.0-or-later

// Package rules implements built-in AIBOM policy rules for AI supply chain
// governance: model pinning, provider restrictions, tool approvals, and more.
package rules

import (
	"fmt"
	"strings"

	"github.com/nostalgicskinco/aibom-policy-engine/pkg/engine"
)

// DenyModel denies a specific model by name.
type DenyModel struct {
	ModelName string
	Reason    string
	Sev       engine.Severity
}

func (r *DenyModel) ID() string          { return "AIBOM-001" }
func (r *DenyModel) Name() string        { return "deny-model" }
func (r *DenyModel) Description() string { return "Deny model: " + r.ModelName }
func (r *DenyModel) Severity() engine.Severity { return r.Sev }

func (r *DenyModel) Evaluate(bom *engine.AIBOM) []engine.Violation {
	var violations []engine.Violation
	for _, c := range bom.Components {
		if c.Type == "model" && (c.Name == r.ModelName || strings.Contains(c.Name, r.ModelName)) {
			violations = append(violations, engine.Violation{
				RuleID:    r.ID(),
				RuleName:  r.Name(),
				Severity:  r.Sev,
				Message:   fmt.Sprintf("model '%s' is denied: %s", c.Name, r.Reason),
				Component: c.Name,
			})
		}
	}
	return violations
}

// RequireModelVersion requires that a specific model has a pinned version.
type RequireModelVersion struct {
	ModelName string
	Sev       engine.Severity
}

func (r *RequireModelVersion) ID() string          { return "AIBOM-002" }
func (r *RequireModelVersion) Name() string        { return "require-model-version" }
func (r *RequireModelVersion) Description() string { return "Model must have pinned version: " + r.ModelName }
func (r *RequireModelVersion) Severity() engine.Severity { return r.Sev }

func (r *RequireModelVersion) Evaluate(bom *engine.AIBOM) []engine.Violation {
	var violations []engine.Violation
	for _, c := range bom.Components {
		if c.Type == "model" && strings.Contains(c.Name, r.ModelName) {
			if c.Version == "" {
				violations = append(violations, engine.Violation{
					RuleID:    r.ID(),
					RuleName:  r.Name(),
					Severity:  r.Sev,
					Message:   fmt.Sprintf("model '%s' has no pinned version", c.Name),
					Component: c.Name,
				})
			}
		}
	}
	return violations
}

// AllowedProviders restricts which providers are allowed.
type AllowedProviders struct {
	Providers []string
	Sev       engine.Severity
}

func (r *AllowedProviders) ID() string          { return "AIBOM-003" }
func (r *AllowedProviders) Name() string        { return "allowed-providers" }
func (r *AllowedProviders) Description() string { return fmt.Sprintf("Only providers: %v", r.Providers) }
func (r *AllowedProviders) Severity() engine.Severity { return r.Sev }

func (r *AllowedProviders) Evaluate(bom *engine.AIBOM) []engine.Violation {
	allowed := make(map[string]bool)
	for _, p := range r.Providers {
		allowed[strings.ToLower(p)] = true
	}

	var violations []engine.Violation
	for _, c := range bom.Components {
		if c.Type == "model" && c.Provider != "" {
			if !allowed[strings.ToLower(c.Provider)] {
				violations = append(violations, engine.Violation{
					RuleID:    r.ID(),
					RuleName:  r.Name(),
					Severity:  r.Sev,
					Message:   fmt.Sprintf("provider '%s' not in allowed list for model '%s'", c.Provider, c.Name),
					Component: c.Name,
				})
			}
		}
	}
	return violations
}

// DenyTool denies a specific tool.
type DenyTool struct {
	ToolName string
	Reason   string
	Sev      engine.Severity
}

func (r *DenyTool) ID() string          { return "AIBOM-004" }
func (r *DenyTool) Name() string        { return "deny-tool" }
func (r *DenyTool) Description() string { return "Deny tool: " + r.ToolName }
func (r *DenyTool) Severity() engine.Severity { return r.Sev }

func (r *DenyTool) Evaluate(bom *engine.AIBOM) []engine.Violation {
	var violations []engine.Violation
	for _, c := range bom.Components {
		if c.Type == "tool" && (c.Name == r.ToolName || strings.Contains(c.Name, r.ToolName)) {
			violations = append(violations, engine.Violation{
				RuleID:    r.ID(),
				RuleName:  r.Name(),
				Severity:  r.Sev,
				Message:   fmt.Sprintf("tool '%s' is denied: %s", c.Name, r.Reason),
				Component: c.Name,
			})
		}
	}
	return violations
}

// MaxModels limits the number of distinct models.
type MaxModels struct {
	Max int
	Sev engine.Severity
}

func (r *MaxModels) ID() string          { return "AIBOM-005" }
func (r *MaxModels) Name() string        { return "max-models" }
func (r *MaxModels) Description() string { return fmt.Sprintf("Max %d models allowed", r.Max) }
func (r *MaxModels) Severity() engine.Severity { return r.Sev }

func (r *MaxModels) Evaluate(bom *engine.AIBOM) []engine.Violation {
	count := 0
	for _, c := range bom.Components {
		if c.Type == "model" {
			count++
		}
	}
	if count > r.Max {
		return []engine.Violation{{
			RuleID:   r.ID(),
			RuleName: r.Name(),
			Severity: r.Sev,
			Message:  fmt.Sprintf("%d models found, max %d allowed", count, r.Max),
		}}
	}
	return nil
}

// MaxTools limits the number of distinct tools.
type MaxTools struct {
	Max int
	Sev engine.Severity
}

func (r *MaxTools) ID() string          { return "AIBOM-006" }
func (r *MaxTools) Name() string        { return "max-tools" }
func (r *MaxTools) Description() string { return fmt.Sprintf("Max %d tools allowed", r.Max) }
func (r *MaxTools) Severity() engine.Severity { return r.Sev }

func (r *MaxTools) Evaluate(bom *engine.AIBOM) []engine.Violation {
	count := 0
	for _, c := range bom.Components {
		if c.Type == "tool" {
			count++
		}
	}
	if count > r.Max {
		return []engine.Violation{{
			RuleID:   r.ID(),
			RuleName: r.Name(),
			Severity: r.Sev,
			Message:  fmt.Sprintf("%d tools found, max %d allowed", count, r.Max),
		}}
	}
	return nil
}

// RequireAllModelsVersioned requires every model to have a version.
type RequireAllModelsVersioned struct {
	Sev engine.Severity
}

func (r *RequireAllModelsVersioned) ID() string          { return "AIBOM-007" }
func (r *RequireAllModelsVersioned) Name() string        { return "require-all-models-versioned" }
func (r *RequireAllModelsVersioned) Description() string { return "All models must have pinned versions" }
func (r *RequireAllModelsVersioned) Severity() engine.Severity { return r.Sev }

func (r *RequireAllModelsVersioned) Evaluate(bom *engine.AIBOM) []engine.Violation {
	var violations []engine.Violation
	for _, c := range bom.Components {
		if c.Type == "model" && c.Version == "" {
			violations = append(violations, engine.Violation{
				RuleID:    r.ID(),
				RuleName:  r.Name(),
				Severity:  r.Sev,
				Message:   fmt.Sprintf("model '%s' has no pinned version", c.Name),
				Component: c.Name,
			})
		}
	}
	return violations
}

// DenyExternalEndpoints denies any external service endpoints.
type DenyExternalEndpoints struct {
	AllowedHosts []string
	Sev          engine.Severity
}

func (r *DenyExternalEndpoints) ID() string          { return "AIBOM-008" }
func (r *DenyExternalEndpoints) Name() string        { return "deny-external-endpoints" }
func (r *DenyExternalEndpoints) Description() string { return "Only allowed external endpoints" }
func (r *DenyExternalEndpoints) Severity() engine.Severity { return r.Sev }

func (r *DenyExternalEndpoints) Evaluate(bom *engine.AIBOM) []engine.Violation {
	allowed := make(map[string]bool)
	for _, h := range r.AllowedHosts {
		allowed[strings.ToLower(h)] = true
	}

	var violations []engine.Violation
	for _, s := range bom.Services {
		host := strings.ToLower(s.Name)
		if !allowed[host] {
			violations = append(violations, engine.Violation{
				RuleID:    r.ID(),
				RuleName:  r.Name(),
				Severity:  r.Sev,
				Message:   fmt.Sprintf("external endpoint '%s' not in allowed list", s.Name),
				Component: s.Name,
			})
		}
	}
	return violations
}
