"""Pytest configuration and fixtures."""
import pytest
from pkg.generator.builder import AIBOMBuilder
from pkg.validator.checker import AIBOMChecker
from pkg.models.aibom import AIBOM, AIComponent, ComponentType

@pytest.fixture
def builder():
    """Create a test builder."""
    return AIBOMBuilder("Test AIBOM", "Test Org")

@pytest.fixture
def checker():
    """Create a test checker."""
    return AIBOMChecker()

@pytest.fixture
def sample_aibom():
    """Create a sample AIBOM."""
    builder = AIBOMBuilder("Sample AIBOM", "Test")
    builder.add_model("GPT-4", "OpenAI", version="1.0")
    builder.add_tool("SearchTool", "Internal")
    return builder.build()
