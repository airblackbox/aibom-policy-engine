// Copyright 2024 Nostalgic Skin Co.
// SPDX-License-Identifier: AGPL-3.0-or-later

package rules

import (
	"testing"

	"github.com/nostalgicskinco/aibom-policy-engine/pkg/engine"
)

func testBOM() *engine.AIBOM {
	return &engine.AIBOM{
		BOMFormat: "CycloneDX-AI",
		Components: []engine.Component{
			{Type: "model", Name: "gpt-4", Version: "2024-01", Provider: "openai"},
			{Type: "model", Name: "claude-3-sonnet", Version: "", Provider: "anthropic"},
			{Type: "tool", Name: "web_search", Version: "1.0"},
			{Type: "tool", Name: "exec_sql"},
			{Type: "framework", Name: "langchain", Version: "0.1.0"},
		},
		Services: []engine.Service{
			{Name: "api.openai.com", Endpoint: "https://api.openai.com/v1"},
			{Name: "api.anthropic.com", Endpoint: "https://api.anthropic.com"},
		},
	}
}

func TestDenyModel(t *testing.T) {
	bom := testBOM()
	rule := &DenyModel{ModelName: "gpt-4", Reason: "not approved", Sev: engine.SeverityHigh}
	violations := rule.Evaluate(bom)
	if len(violations) == 0 {
		t.Fatal("expected violation for gpt-4")
	}

	rule2 := &DenyModel{ModelName: "llama-3", Reason: "test", Sev: engine.SeverityHigh}
	violations2 := rule2.Evaluate(bom)
	if len(violations2) != 0 {
		t.Fatal("expected no violation for llama-3")
	}
}

func TestRequireModelVersion(t *testing.T) {
	bom := testBOM()
	// claude-3-sonnet has no version
	rule := &RequireModelVersion{ModelName: "claude-3-sonnet", Sev: engine.SeverityHigh}
	violations := rule.Evaluate(bom)
	if len(violations) == 0 {
		t.Fatal("expected violation for unversioned claude-3-sonnet")
	}

	// gpt-4 has version
	rule2 := &RequireModelVersion{ModelName: "gpt-4", Sev: engine.SeverityHigh}
	violations2 := rule2.Evaluate(bom)
	if len(violations2) != 0 {
		t.Fatal("expected no violation for versioned gpt-4")
	}
}

func TestAllowedProviders(t *testing.T) {
	bom := testBOM()
	// Only openai allowed
	rule := &AllowedProviders{Providers: []string{"openai"}, Sev: engine.SeverityHigh}
	violations := rule.Evaluate(bom)
	if len(violations) == 0 {
		t.Fatal("expected violation for anthropic provider")
	}

	// Both allowed
	rule2 := &AllowedProviders{Providers: []string{"openai", "anthropic"}, Sev: engine.SeverityHigh}
	violations2 := rule2.Evaluate(bom)
	if len(violations2) != 0 {
		t.Fatal("expected no violations with both providers allowed")
	}
}

func TestDenyTool(t *testing.T) {
	bom := testBOM()
	rule := &DenyTool{ToolName: "exec_sql", Reason: "dangerous", Sev: engine.SeverityCritical}
	violations := rule.Evaluate(bom)
	if len(violations) == 0 {
		t.Fatal("expected violation for exec_sql")
	}
	if violations[0].Severity != engine.SeverityCritical {
		t.Fatalf("expected critical severity, got %s", violations[0].Severity)
	}
}

func TestMaxModels(t *testing.T) {
	bom := testBOM()
	// 2 models, max 1 → violate
	rule := &MaxModels{Max: 1, Sev: engine.SeverityMedium}
	violations := rule.Evaluate(bom)
	if len(violations) == 0 {
		t.Fatal("expected violation for > 1 model")
	}

	// max 5 → pass
	rule2 := &MaxModels{Max: 5, Sev: engine.SeverityMedium}
	violations2 := rule2.Evaluate(bom)
	if len(violations2) != 0 {
		t.Fatal("expected no violation for max 5")
	}
}

func TestMaxTools(t *testing.T) {
	bom := testBOM()
	rule := &MaxTools{Max: 1, Sev: engine.SeverityMedium}
	violations := rule.Evaluate(bom)
	if len(violations) == 0 {
		t.Fatal("expected violation for > 1 tool")
	}
}

func TestRequireAllModelsVersioned(t *testing.T) {
	bom := testBOM()
	rule := &RequireAllModelsVersioned{Sev: engine.SeverityHigh}
	violations := rule.Evaluate(bom)
	// claude-3-sonnet has no version
	if len(violations) == 0 {
		t.Fatal("expected violation for unversioned model")
	}
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestDenyExternalEndpoints(t *testing.T) {
	bom := testBOM()
	// Only allow openai
	rule := &DenyExternalEndpoints{
		AllowedHosts: []string{"api.openai.com"},
		Sev:          engine.SeverityHigh,
	}
	violations := rule.Evaluate(bom)
	if len(violations) == 0 {
		t.Fatal("expected violation for api.anthropic.com")
	}

	// Allow both
	rule2 := &DenyExternalEndpoints{
		AllowedHosts: []string{"api.openai.com", "api.anthropic.com"},
		Sev:          engine.SeverityHigh,
	}
	violations2 := rule2.Evaluate(bom)
	if len(violations2) != 0 {
		t.Fatal("expected no violations with both hosts allowed")
	}
}

func TestEngineIntegration(t *testing.T) {
	bom := testBOM()
	eng := engine.NewEngine()
	eng.AddPolicy(engine.Policy{
		Name: "governance",
		Rules: []engine.Rule{
			&DenyTool{ToolName: "exec_sql", Reason: "forbidden", Sev: engine.SeverityCritical},
			&AllowedProviders{Providers: []string{"openai", "anthropic"}, Sev: engine.SeverityHigh},
			&RequireAllModelsVersioned{Sev: engine.SeverityMedium},
		},
	})

	results := eng.EvaluateAll(bom)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := results[0]
	if r.CriticalCount != 1 {
		t.Fatalf("expected 1 critical, got %d", r.CriticalCount)
	}
	if r.Passed() {
		t.Fatal("expected policy to fail")
	}
}

func TestRegistryLoadJSON(t *testing.T) {
	data := []byte(`{
		"name": "test-policy",
		"rules": [
			{"type": "deny-model", "severity": "critical", "params": {"model_name": "gpt-4", "reason": "not approved"}},
			{"type": "max-models", "severity": "medium", "params": {"max": 3}},
			{"type": "require-all-models-versioned", "severity": "high"}
		]
	}`)

	registry := DefaultRegistry()
	pol, err := registry.ParsePolicyJSON(data)
	if err != nil {
		t.Fatalf("ParsePolicyJSON: %v", err)
	}
	if len(pol.Rules) != 3 {
		t.Fatalf("expected 3 rules, got %d", len(pol.Rules))
	}

	bom := testBOM()
	eng := engine.NewEngine()
	eng.AddPolicy(*pol)
	results := eng.EvaluateAll(bom)
	if !engine.HasFailures(results) {
		t.Fatal("expected failures from deny-model")
	}
}
