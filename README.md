# aibom-policy-engine

**AIBOM policy-as-code engine** — deny risky combinations of models, tools, and data in your AI supply chain.

OPA-style declarative policies for AI governance: model pinning, provider restrictions, tool denylists, endpoint allowlists, and component count limits. Evaluates against CycloneDX-AI Bills of Materials.

> Part of the **GenAI Infrastructure Standard** — a composable suite of open-source tools for enterprise-grade GenAI observability, security, and governance.
>
> | Layer | Component | Repo |
> |-------|-----------|------|
> | Privacy | Prompt Vault Processor | [prompt-vault-processor](https://github.com/nostalgicskinco/prompt-vault-processor) |
> | Normalization | Semantic Normalizer | [genai-semantic-normalizer](https://github.com/nostalgicskinco/genai-semantic-normalizer) |
> | Metrics | Cost & SLO Pack | [genai-cost-slo](https://github.com/nostalgicskinco/genai-cost-slo) |
> | Replay | Agent VCR | [agent-vcr](https://github.com/nostalgicskinco/agent-vcr) |
> | Testing | Regression Harness | [trace-regression-harness](https://github.com/nostalgicskinco/trace-regression-harness) |
> | Security | MCP Scanner | [mcp-security-scanner](https://github.com/nostalgicskinco/mcp-security-scanner) |
> | Gateway | MCP Policy Gateway | [mcp-policy-gateway](https://github.com/nostalgicskinco/mcp-policy-gateway) |
> | Inventory | Runtime AIBOM Emitter | [runtime-aibom-emitter](https://github.com/nostalgicskinco/runtime-aibom-emitter) |
> | **Policy** | **AIBOM Policy Engine** | **this repo** |

## Quick Start

```bash
go build -o aibompolicy ./cmd/aibompolicy
./aibompolicy -bom aibom.json -policy governance.json
```

## Built-in Rules

| Rule | ID | Description |
|------|----|-------------|
| `deny-model` | AIBOM-001 | Deny specific models by name |
| `require-model-version` | AIBOM-002 | Require pinned version for a model |
| `allowed-providers` | AIBOM-003 | Restrict to approved providers only |
| `deny-tool` | AIBOM-004 | Deny specific tools |
| `max-models` | AIBOM-005 | Limit distinct model count |
| `max-tools` | AIBOM-006 | Limit distinct tool count |
| `require-all-models-versioned` | AIBOM-007 | All models must have versions |
| `deny-external-endpoints` | AIBOM-008 | Only allow approved endpoints |

## License

AGPL-3.0-or-later — see [LICENSE](LICENSE). Commercial licenses available.
