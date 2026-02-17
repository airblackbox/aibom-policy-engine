"""Test FastAPI routes."""
import pytest
from fastapi.testclient import TestClient
from pkg.api.routes import router

client = TestClient(router)

def test_health():
    """Test health endpoint."""
    response = client.get("/v1/health")
    assert response.status_code == 200
    data = response.json()
    assert data["status"] == "ok"

def test_create_aibom():
    """Test creating AIBOM."""
    response = client.post(
        "/v1/aibom/create",
        json={
            "name": "Test AIBOM",
            "organization": "TestOrg",
            "components": []
        }
    )
    assert response.status_code == 200
    aibom = response.json()
    assert aibom["name"] == "Test AIBOM"
    assert "id" in aibom

def test_get_aibom():
    """Test getting AIBOM."""
    create_resp = client.post(
        "/v1/aibom/create",
        json={"name": "Test", "organization": ""}
    )
    aibom_id = create_resp.json()["id"]
    get_resp = client.get(f"/v1/aibom/{aibom_id}")
    assert get_resp.status_code == 200
    assert get_resp.json()["id"] == aibom_id

def test_get_nonexistent_aibom():
    """Test getting nonexistent AIBOM."""
    response = client.get("/v1/aibom/nonexistent")
    assert response.status_code == 404

def test_validate_aibom():
    """Test validating AIBOM."""
    create_resp = client.post(
        "/v1/aibom/create",
        json={"name": "Test", "organization": ""}
    )
    aibom_id = create_resp.json()["id"]
    validate_resp = client.post(f"/v1/aibom/{aibom_id}/validate")
    assert validate_resp.status_code == 200
    result = validate_resp.json()
    assert "valid" in result

def test_list_aiboms():
    """Test listing AIBOMs."""
    response = client.get("/v1/aiboms")
    assert response.status_code == 200
    data = response.json()
    assert "count" in data
    assert "aiboms" in data
