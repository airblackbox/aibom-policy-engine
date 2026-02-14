// Copyright 2024 Nostalgic Skin Co.
// SPDX-License-Identifier: AGPL-3.0-or-later

// Package engine evaluates AIBOM policy rules against AI Bill of Materials.
// Policies express constraints on which models, tools, providers, and
// configurations are permitted in an AI system.
package engine

import (
	"encoding/json"
	"fmt"
	"os"
)

// Severity indicates the impact level of a violation.
type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"
)

// AIBOM is the input format — a simplified CycloneDX-AI BOM.
type AIBOM struct {
	BOMFormat   string      `json:"bomFormat"`
	Components  []Component `json:"components"`
	Services    []Service   `json:"services,omitempty"`
}

// Component is a model, tool, library, or framework in the BOM.
type Component struct {
	Type       string            `json:"type"`
	Name       string            `json:"name"`
	Version    string            `json:"version,omitempty"`
	Provider   string            `json:"provider,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

// Service is an external endpoint in the BOM.
type Service struct {
	Name     string `json:"name"`
	Endpoint string `json:"endpoint,omitempty"`
	Provider string `json:"provider,omitempty"`
}

// Rule is a single policy assertion evaluated against an AIBOM.
type Rule interface {
	ID() string
	Name() string
	Description() string
	Severity() Severity
	Evaluate(bom *AIBOM) []Violation
}

// Violation is a policy rule failure.
type Violation struct {
	RuleID      string   `json:"ruleId"`
	RuleName    string   `json:"ruleName"`
	Severity    Severity `json:"severity"`
	Message     string   `json:"message"`
	Component   string   `json:"component,omitempty"`
}

// PolicyFile is a JSON-serializable policy configuration.
type PolicyFile struct {
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Rules       []RuleDef `json:"rules"`
}

// RuleDef defines a rule in a policy file.
type RuleDef struct {
	Type     string         `json:"type"`
	Severity string         `json:"severity,omitempty"`
	Params   map[string]any `json:"params,omitempty"`
}

// Result holds evaluation results for one policy.
type Result struct {
	PolicyName    string      `json:"policyName"`
	Violations    []Violation `json:"violations"`
	CriticalCount int        `json:"criticalCount"`
	HighCount     int        `json:"highCount"`
	MediumCount   int        `json:"mediumCount"`
	LowCount      int        `json:"lowCount"`
}

// Passed returns true if no critical or high violations exist.
func (r *Result) Passed() bool {
	return r.CriticalCount == 0 && r.HighCount == 0
}

// Engine evaluates policies against AIBOMs.
type Engine struct {
	policies []Policy
}

// Policy is a named collection of rules.
type Policy struct {
	Name  string
	Rules []Rule
}

// NewEngine creates a new policy engine.
func NewEngine() *Engine {
	return &Engine{}
}

// AddPolicy adds a policy.
func (e *Engine) AddPolicy(p Policy) {
	e.policies = append(e.policies, p)
}

// EvaluateAll runs all policies against the given AIBOM.
func (e *Engine) EvaluateAll(bom *AIBOM) []*Result {
	var results []*Result
	for _, p := range e.policies {
		r := &Result{PolicyName: p.Name}
		for _, rule := range p.Rules {
			for _, v := range rule.Evaluate(bom) {
				r.Violations = append(r.Violations, v)
				switch v.Severity {
				case SeverityCritical:
					r.CriticalCount++
				case SeverityHigh:
					r.HighCount++
				case SeverityMedium:
					r.MediumCount++
				case SeverityLow:
					r.LowCount++
				}
			}
		}
		results = append(results, r)
	}
	return results
}

// HasFailures returns true if any result has critical or high violations.
func HasFailures(results []*Result) bool {
	for _, r := range results {
		if !r.Passed() {
			return true
		}
	}
	return false
}

// LoadBOM reads an AIBOM from a JSON file.
func LoadBOM(path string) (*AIBOM, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read BOM: %w", err)
	}
	var b AIBOM
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("parse BOM: %w", err)
	}
	return &b, nil
}
