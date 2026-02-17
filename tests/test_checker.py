"""Test AIBOMChecker."""
import pytest
from pkg.validator.checker import AIBOMChecker
from pkg.models.aibom import AIBOM, AIComponent, ComponentType, RiskClassification

@pytest.fixture
def checker():
    """Create checker."""
    return AIBOMChecker()

def test_validate_empty_aibom(checker):
    """Test validating empty AIBOM."""
    aibom = AIBOM(name="Empty")
    result = checker.validate(aibom)
    assert result.valid is True

def test_validate_missing_component_id(checker):
    """Test validation catches missing IDs."""
    aibom = AIBOM(name="Test")
    comp = AIComponent(id="", name="NoID", component_type=ComponentType.MODEL)
    aibom.components = [comp]
    result = checker.validate(aibom)
    assert result.valid is False
    assert any("missing ID" in e for e in result.errors)

def test_validate_duplicate_ids(checker):
    """Test validation catches duplicate IDs."""
    aibom = AIBOM(name="Test")
    comp1 = AIComponent(id="same", name="c1", component_type=ComponentType.MODEL)
    comp2 = AIComponent(id="same", name="c2", component_type=ComponentType.TOOL)
    aibom.components = [comp1, comp2]
    result = checker.validate(aibom)
    assert result.valid is False
    assert any("Duplicate" in e for e in result.errors)

def test_validate_invalid_dependency(checker):
    """Test validation catches invalid dependencies."""
    aibom = AIBOM(name="Test")
    comp = AIComponent(id="c1", name="Model", component_type=ComponentType.MODEL)
    aibom.components = [comp]
    aibom.dependencies = [{"from": "c1", "to": "unknown"}]
    result = checker.validate(aibom)
    assert result.valid is False

def test_validate_high_risk_no_description(checker):
    """Test warning for high-risk without description."""
    aibom = AIBOM(name="Test")
    comp = AIComponent(
        id="c1",
        name="Risky",
        component_type=ComponentType.MODEL,
        risk_classification=RiskClassification.HIGH,
        description=""
    )
    aibom.components = [comp]
    result = checker.validate(aibom)
    assert result.valid is True
    assert any("missing description" in w for w in result.warnings)

def test_validate_model_no_provider(checker):
    """Test warning for model without provider."""
    aibom = AIBOM(name="Test")
    comp = AIComponent(
        id="c1",
        name="Model",
        component_type=ComponentType.MODEL,
        provider=""
    )
    aibom.components = [comp]
    result = checker.validate(aibom)
    assert result.valid is True
    assert any("missing provider" in w for w in result.warnings)
