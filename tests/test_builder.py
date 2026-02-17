"""Test AIBOMBuilder."""
import pytest
from pkg.generator.builder import AIBOMBuilder
from pkg.models.aibom import ComponentType, RiskClassification

def test_builder_initialization():
    """Test builder initialization."""
    builder = AIBOMBuilder("Test", "TestOrg")
    assert builder.name == "Test"
    assert builder.organization == "TestOrg"

def test_add_model():
    """Test adding a model."""
    builder = AIBOMBuilder("Test")
    comp_id = builder.add_model("GPT-4", "OpenAI", version="1.0")
    assert comp_id.startswith("model-")
    aibom = builder.build()
    assert len(aibom.components) == 1
    assert aibom.components[0].name == "GPT-4"

def test_add_tool():
    """Test adding a tool."""
    builder = AIBOMBuilder("Test")
    comp_id = builder.add_tool("SearchTool", "Internal")
    assert comp_id.startswith("tool-")
    aibom = builder.build()
    assert len(aibom.components) == 1
    assert aibom.components[0].name == "SearchTool"

def test_add_data_source():
    """Test adding a data source."""
    builder = AIBOMBuilder("Test")
    comp_id = builder.add_data_source("UserDB", "Internal")
    assert comp_id.startswith("data-")
    aibom = builder.build()
    assert aibom.components[0].component_type == ComponentType.DATA_SOURCE

def test_add_dependency():
    """Test adding dependency."""
    builder = AIBOMBuilder("Test")
    id1 = builder.add_model("Model1", "Provider1")
    id2 = builder.add_tool("Tool1")
    builder.add_dependency(id1, id2)
    aibom = builder.build()
    assert len(aibom.dependencies) == 1

def test_set_metadata():
    """Test setting metadata."""
    builder = AIBOMBuilder("Test")
    builder.set_metadata("version_info", "v1.0.0")
    aibom = builder.build()
    assert aibom.metadata["version_info"] == "v1.0.0"

def test_build():
    """Test building AIBOM."""
    builder = AIBOMBuilder("Test AIBOM", "TestOrg")
    builder.add_model("GPT-4", "OpenAI")
    builder.add_tool("SearchTool")
    aibom = builder.build()
    assert aibom.name == "Test AIBOM"
    assert aibom.organization == "TestOrg"
    assert len(aibom.components) == 2
    assert aibom.id.startswith("aibom-")

def test_from_gateway():
    """Test building from gateway data."""
    gateway_data = {
        "name": "Agent System",
        "organization": "Test",
        "models": [
            {"name": "GPT-4", "provider": "OpenAI", "version": "1.0"}
        ],
        "tools": [
            {"name": "SearchTool", "provider": "Internal"}
        ],
    }
    aibom = AIBOMBuilder.from_gateway(gateway_data)
    assert aibom.name == "Agent System"
    assert len(aibom.components) == 2
