"""Test AIBOM models."""
import pytest
from pkg.models.aibom import (
    ComponentType,
    RiskClassification,
    AIComponent,
    AIBOM,
    AIBOMValidation,
)

def test_component_type_enum():
    """Test ComponentType enum."""
    assert ComponentType.MODEL.value == "model"
    assert ComponentType.TOOL.value == "tool"
    assert ComponentType.DATA_SOURCE.value == "data_source"

def test_risk_classification_enum():
    """Test RiskClassification enum."""
    assert RiskClassification.MINIMAL.value == "minimal"
    assert RiskClassification.HIGH.value == "high"

def test_ai_component_creation():
    """Test AIComponent creation."""
    comp = AIComponent(
        id="test-1",
        name="TestModel",
        component_type=ComponentType.MODEL,
        provider="TestProvider"
    )
    assert comp.id == "test-1"
    assert comp.name == "TestModel"
    assert comp.provider == "TestProvider"

def test_aibom_creation():
    """Test AIBOM creation."""
    aibom = AIBOM(name="Test AIBOM")
    assert aibom.name == "Test AIBOM"
    assert aibom.version == "1.0"
    assert len(aibom.components) == 0

def test_aibom_model_count():
    """Test model count property."""
    aibom = AIBOM(name="Test")
    comp1 = AIComponent(id="1", name="m1", component_type=ComponentType.MODEL)
    comp2 = AIComponent(id="2", name="m2", component_type=ComponentType.MODEL)
    comp3 = AIComponent(id="3", name="t1", component_type=ComponentType.TOOL)
    aibom.components = [comp1, comp2, comp3]
    assert aibom.model_count == 2

def test_aibom_high_risk_components():
    """Test high risk components property."""
    aibom = AIBOM(name="Test")
    comp1 = AIComponent(
        id="1",
        name="safe",
        component_type=ComponentType.MODEL,
        risk_classification=RiskClassification.MINIMAL
    )
    comp2 = AIComponent(
        id="2",
        name="risky",
        component_type=ComponentType.MODEL,
        risk_classification=RiskClassification.HIGH
    )
    aibom.components = [comp1, comp2]
    assert len(aibom.high_risk_components) == 1
    assert aibom.high_risk_components[0].name == "risky"

def test_aibom_validation():
    """Test AIBOMValidation model."""
    validation = AIBOMValidation(valid=True)
    assert validation.valid is True
    assert len(validation.errors) == 0
