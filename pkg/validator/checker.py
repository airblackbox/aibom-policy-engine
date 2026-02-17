"""AIBOM validation checker."""
from __future__ import annotations
from pkg.models.aibom import AIBOM, AIBOMValidation, RiskClassification

class AIBOMChecker:
    """Validates AIBOM documents."""
    def validate(self, aibom: AIBOM) -> AIBOMValidation:
        """Validate an AIBOM document."""
        errors = []
        warnings = []
        
        # Check all components have IDs
        for i, comp in enumerate(aibom.components):
            if not comp.id:
                errors.append(f"Component {i} missing ID")
        
        # Check for duplicate IDs
        ids = [c.id for c in aibom.components if c.id]
        if len(ids) != len(set(ids)):
            errors.append("Duplicate component IDs found")
        
        # Check dependencies reference valid components
        valid_ids = {c.id for c in aibom.components}
        for dep in aibom.dependencies:
            from_id = dep.get("from")
            to_id = dep.get("to")
            if from_id not in valid_ids:
                errors.append(f"Dependency references unknown component: {from_id}")
            if to_id not in valid_ids:
                errors.append(f"Dependency references unknown component: {to_id}")
        
        # Check high-risk components have descriptions
        for comp in aibom.high_risk_components:
            if not comp.description:
                warnings.append(
                    f"High-risk component '{comp.name}' missing description"
                )
        
        # Check models have providers
        for comp in aibom.components:
            if comp.component_type.value == "model" and not comp.provider:
                warnings.append(f"Model '{comp.name}' missing provider")
        
        valid = len(errors) == 0
        return AIBOMValidation(
            valid=valid,
            errors=errors,
            warnings=warnings,
        )
