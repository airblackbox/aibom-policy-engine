"""AIBOM data models."""
from __future__ import annotations
from datetime import datetime
from enum import Enum
from typing import Any
from pydantic import BaseModel, Field

class ComponentType(str, Enum):
    """AI component types."""
    MODEL = "model"
    TOOL = "tool"
    DATA_SOURCE = "data_source"
    POLICY = "policy"
    PROCESSOR = "processor"
    FRAMEWORK = "framework"

class RiskClassification(str, Enum):
    """Risk classification levels."""
    MINIMAL = "minimal"
    LIMITED = "limited"
    HIGH = "high"
    UNACCEPTABLE = "unacceptable"

class AIComponent(BaseModel):
    """Individual AI component in AIBOM."""
    id: str
    name: str
    version: str = ""
    component_type: ComponentType
    provider: str = ""
    risk_classification: RiskClassification = RiskClassification.MINIMAL
    description: str = ""
    license: str = ""
    capabilities: list[str] = Field(default_factory=list)
    limitations: list[str] = Field(default_factory=list)
    metadata: dict[str, Any] = Field(default_factory=dict)

class AIBOM(BaseModel):
    """AI Bill of Materials document."""
    id: str = ""
    name: str
    version: str = "1.0"
    created_at: datetime = Field(default_factory=datetime.utcnow)
    organization: str = ""
    components: list[AIComponent] = Field(default_factory=list)
    dependencies: list[dict[str, str]] = Field(default_factory=list)
    metadata: dict[str, Any] = Field(default_factory=dict)

    @property
    def model_count(self) -> int:
        """Count models."""
        return sum(
            1 for c in self.components
            if c.component_type == ComponentType.MODEL
        )

    @property
    def tool_count(self) -> int:
        """Count tools."""
        return sum(
            1 for c in self.components
            if c.component_type == ComponentType.TOOL
        )

    @property
    def high_risk_components(self) -> list[AIComponent]:
        """Get high-risk components."""
        return [
            c for c in self.components
            if c.risk_classification in (
                RiskClassification.HIGH,
                RiskClassification.UNACCEPTABLE
            )
        ]

class AIBOMValidation(BaseModel):
    """AIBOM validation result."""
    valid: bool = True
    errors: list[str] = Field(default_factory=list)
    warnings: list[str] = Field(default_factory=list)
