"""AIBOM builder for generation."""
from __future__ import annotations
import uuid
from typing import Any
from pkg.models.aibom import (
    AIComponent,
    AIBOM,
    ComponentType,
    RiskClassification,
)

class AIBOMBuilder:
    """Builder for constructing AIBOM documents."""
    def __init__(self, name: str, organization: str = "") -> None:
        self.name = name
        self.organization = organization
        self._components: list[AIComponent] = []
        self._dependencies: list[dict[str, str]] = []
        self._metadata: dict[str, Any] = {}

    def add_model(
        self,
        name: str,
        provider: str,
        version: str = "",
        risk: RiskClassification = RiskClassification.MINIMAL,
        description: str = "",
    ) -> str:
        """Add a model component."""
        comp_id = f"model-{uuid.uuid4().hex[:8]}"
        component = AIComponent(
            id=comp_id,
            name=name,
            version=version,
            component_type=ComponentType.MODEL,
            provider=provider,
            risk_classification=risk,
            description=description,
        )
        self._components.append(component)
        return comp_id

    def add_tool(
        self,
        name: str,
        provider: str = "",
        version: str = "",
        description: str = "",
    ) -> str:
        """Add a tool component."""
        comp_id = f"tool-{uuid.uuid4().hex[:8]}"
        component = AIComponent(
            id=comp_id,
            name=name,
            version=version,
            component_type=ComponentType.TOOL,
            provider=provider,
            description=description,
        )
        self._components.append(component)
        return comp_id

    def add_data_source(
        self,
        name: str,
        provider: str = "",
        description: str = "",
    ) -> str:
        """Add a data source component."""
        comp_id = f"data-{uuid.uuid4().hex[:8]}"
        component = AIComponent(
            id=comp_id,
            name=name,
            component_type=ComponentType.DATA_SOURCE,
            provider=provider,
            description=description,
        )
        self._components.append(component)
        return comp_id

    def add_dependency(self, from_id: str, to_id: str) -> None:
        """Add dependency between components."""
        self._dependencies.append({"from": from_id, "to": to_id})

    def set_metadata(self, key: str, value: Any) -> None:
        """Set metadata field."""
        self._metadata[key] = value

    def build(self) -> AIBOM:
        """Build AIBOM document."""
        aibom = AIBOM(
            id=f"aibom-{uuid.uuid4().hex[:12]}",
            name=self.name,
            organization=self.organization,
            components=self._components,
            dependencies=self._dependencies,
            metadata=self._metadata,
        )
        return aibom

    @staticmethod
    def from_gateway(gateway_data: dict[str, Any]) -> AIBOM:
        """Build AIBOM from gateway health data."""
        builder = AIBOMBuilder(
            name=gateway_data.get("name", "Generated AIBOM"),
            organization=gateway_data.get("organization", ""),
        )
        models = gateway_data.get("models", [])
        for model_info in models:
            builder.add_model(
                name=model_info.get("name"),
                provider=model_info.get("provider"),
                version=model_info.get("version", ""),
                description=model_info.get("description", ""),
            )
        tools = gateway_data.get("tools", [])
        for tool_info in tools:
            builder.add_tool(
                name=tool_info.get("name"),
                provider=tool_info.get("provider", ""),
                description=tool_info.get("description", ""),
            )
        return builder.build()
