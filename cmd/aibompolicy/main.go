// Copyright 2024 Nostalgic Skin Co.
// SPDX-License-Identifier: AGPL-3.0-or-later

// Command aibompolicy evaluates AIBOM policy rules against an AI BOM file.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/nostalgicskinco/aibom-policy-engine/pkg/engine"
	"github.com/nostalgicskinco/aibom-policy-engine/pkg/rules"
)

func main() {
	bomFile := flag.String("bom", "", "Path to AIBOM JSON file")
	policyFile := flag.String("policy", "", "Path to policy JSON file")
	format := flag.String("format", "text", "Output format: text or json")
	flag.Parse()

	if *bomFile == "" || *policyFile == "" {
		fmt.Fprintf(os.Stderr, "Usage: aibompolicy -bom <aibom.json> -policy <policy.json>\n")
		os.Exit(1)
	}

	bom, err := engine.LoadBOM(*bomFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading BOM: %v\n", err)
		os.Exit(1)
	}

	registry := rules.DefaultRegistry()
	pol, err := registry.LoadPolicyFile(*policyFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading policy: %v\n", err)
		os.Exit(1)
	}

	eng := engine.NewEngine()
	eng.AddPolicy(*pol)
	results := eng.EvaluateAll(bom)

	switch *format {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(results)
	default:
		for _, r := range results {
			fmt.Printf("Policy: %s\n", r.PolicyName)
			if len(r.Violations) == 0 {
				fmt.Println("  ✓ All checks passed")
				continue
			}
			for _, v := range r.Violations {
				fmt.Printf("  ✗ [%s] %s: %s\n", strings.ToUpper(string(v.Severity)), v.RuleID, v.Message)
			}
			fmt.Printf("  Summary: %d critical, %d high, %d medium, %d low\n",
				r.CriticalCount, r.HighCount, r.MediumCount, r.LowCount)
		}
	}

	if engine.HasFailures(results) {
		os.Exit(2)
	}
}
