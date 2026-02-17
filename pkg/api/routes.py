"""FastAPI routes for AIBOM."""
from __future__ import annotations
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from pkg.models.aibom import AIBOM, AIComponent, AIBOMValidation
from pkg.generator.builder import AIBOMBuilder
from pkg.validator.checker import AIBOMChecker

router = FastAPI(title="AIBOM Policy Engine")
checker = AIBOMChecker()
_aiboms: dict[str, AIBOM] = {}

class ComponentInput(BaseModel):
    """Component input model."""
    name: str
    component_type: str
    provider: str = ""
    version: str = ""
    description: str = ""
    risk_classification: str = "minimal"

class AIBOMInput(BaseModel):
    """AIBOM creation input."""
    name: str
    organization: str = ""
    components: list[ComponentInput] = []

@router.get("/v1/health")
async def health():
    """Health check endpoint."""
    return {
        "status": "ok",
        "service": "aibom-policy-engine",
        "aiboms_stored": len(_aiboms)
    }

@router.post("/v1/aibom/create")
async def create_aibom(input_data: AIBOMInput) -> AIBOM:
    """Create a new AIBOM."""
    builder = AIBOMBuilder(input_data.name, input_data.organization)
    for comp_input in input_data.components:
        if comp_input.component_type == "model":
            builder.add_model(
                name=comp_input.name,
                provider=comp_input.provider,
                version=comp_input.version,
                description=comp_input.description,
            )
        elif comp_input.component_type == "tool":
            builder.add_tool(
                name=comp_input.name,
                provider=comp_input.provider,
                version=comp_input.version,
                description=comp_input.description,
            )
    aibom = builder.build()
    _aiboms[aibom.id] = aibom
    return aibom

@router.get("/v1/aibom/{aibom_id}")
async def get_aibom(aibom_id: str) -> AIBOM:
    """Get AIBOM by ID."""
    if aibom_id not in _aiboms:
        raise HTTPException(status_code=404, detail="AIBOM not found")
    return _aiboms[aibom_id]

@router.post("/v1/aibom/{aibom_id}/validate")
async def validate_aibom(aibom_id: str) -> AIBOMValidation:
    """Validate an AIBOM."""
    if aibom_id not in _aiboms:
        raise HTTPException(status_code=404, detail="AIBOM not found")
    aibom = _aiboms[aibom_id]
    return checker.validate(aibom)

@router.post("/v1/components")
async def add_component(aibom_id: str, component: ComponentInput):
    """Add component to AIBOM."""
    if aibom_id not in _aiboms:
        raise HTTPException(status_code=404, detail="AIBOM not found")
    aibom = _aiboms[aibom_id]
    comp = AIComponent(
        id=f"comp-{len(aibom.components)}",
        name=component.name,
        component_type=component.component_type,
        provider=component.provider,
        version=component.version,
        description=component.description,
    )
    aibom.components.append(comp)
    return {"added": True, "component_id": comp.id}

@router.get("/v1/aiboms")
async def list_aiboms():
    """List all AIBOMs."""
    return {
        "count": len(_aiboms),
        "aiboms": [{"id": id, "name": aibom.name} 
                   for id, aibom in _aiboms.items()]
    }
