# AIBOM Policy Engine

[![CI](https://github.com/airblackbox/aibom-policy-engine/actions/workflows/ci.yml/badge.svg)](https://github.com/airblackbox/aibom-policy-engine/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://github.com/airblackbox/aibom-policy-engine/blob/main/LICENSE)
[![Python 3.10+](https://img.shields.io/badge/python-3.10+-3776AB.svg?logo=python&logoColor=white)](https://python.org)


AI Bill of Materials (AIBOM) generation and validation engine. Creates structured inventories of all AI components in a system (models, tools, data sources, policies).

## Features

- Build AIBOM from components or gateway data
- Validate AIBOM structure and content
- EU AI Act risk classifications
- Dependency tracking
- RESTful API with FastAPI
- CLI for AIBOM management

## Quick Start

```bash
pip install -e .
python -m app.server
```

API runs on `http://localhost:8600/v1`

## AIBOM Components

- **Models**: ML/LLM models (GPT-4, Claude, etc.)
- **Tools**: Callable functions/APIs
- **Data Sources**: Databases, APIs, document stores
- **Policies**: Rules and constraints
- **Processors**: Data transformers
- **Frameworks**: Libraries (PyTorch, TensorFlow)

## API Endpoints

| Method | Endpoint | Purpose |
|--------|----------|---------|
| GET | `/v1/health` | Health check |
| POST | `/v1/aibom/create` | Create AIBOM |
| GET | `/v1/aibom/{id}` | Get AIBOM |
| POST | `/v1/aibom/{id}/validate` | Validate AIBOM |
| POST | `/v1/components` | Add component |
| GET | `/v1/aiboms` | List all AIBOMs |

## Risk Classifications

- MINIMAL: Low risk components
- LIMITED: Moderate risk
- HIGH: Significant risk
- UNACCEPTABLE: Should not be used

## Testing

```bash
pytest tests/ -v
```

## EU AI Act Reference

AIBOM aligns with EU AI Act transparency requirements for system components and data sources.

## License

MIT
