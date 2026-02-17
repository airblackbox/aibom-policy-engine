"""CLI commands for AIBOM."""
import click
from rich.console import Console
from rich.table import Table
import httpx
import json

console = Console()
BASE_URL = "http://localhost:8600/v1"

@click.group()
def cli():
    """AIBOM Policy Engine CLI."""
    pass

@cli.command()
def health():
    """Check service health."""
    try:
        with httpx.Client() as client:
            resp = client.get(f"{BASE_URL}/health")
            resp.raise_for_status()
            data = resp.json()
            console.print("[green]✓[/green] Service is healthy")
            console.print(f"  AIBOMs stored: {data['aiboms_stored']}")
    except Exception as e:
        console.print(f"[red]✗[/red] Health check failed: {e}")

@cli.command()
@click.option("--name", required=True, help="AIBOM name")
@click.option("--org", default="", help="Organization")
def create(name: str, org: str):
    """Create a new AIBOM."""
    try:
        with httpx.Client() as client:
            resp = client.post(
                f"{BASE_URL}/aibom/create",
                json={"name": name, "organization": org}
            )
            resp.raise_for_status()
            aibom = resp.json()
            console.print(f"[green]✓[/green] Created AIBOM: {aibom['id']}")
            console.print(f"  Name: {aibom['name']}")
    except Exception as e:
        console.print(f"[red]Error:[/red] {e}")

@cli.command()
def list_aiboms():
    """List all AIBOMs."""
    try:
        with httpx.Client() as client:
            resp = client.get(f"{BASE_URL}/aiboms")
            resp.raise_for_status()
            data = resp.json()
            if data["aiboms"]:
                table = Table(title="AIBOMs")
                table.add_column("ID")
                table.add_column("Name")
                for item in data["aiboms"]:
                    table.add_row(item["id"], item["name"])
                console.print(table)
            else:
                console.print("No AIBOMs found")
    except Exception as e:
        console.print(f"[red]Error:[/red] {e}")

@cli.command()
@click.option("--id", "aibom_id", required=True, help="AIBOM ID")
def validate(aibom_id: str):
    """Validate an AIBOM."""
    try:
        with httpx.Client() as client:
            resp = client.post(f"{BASE_URL}/aibom/{aibom_id}/validate")
            resp.raise_for_status()
            result = resp.json()
            if result["valid"]:
                console.print("[green]✓[/green] AIBOM is valid")
            else:
                console.print("[red]✗[/red] AIBOM has errors:")
                for error in result["errors"]:
                    console.print(f"  - {error}")
            if result["warnings"]:
                console.print("[yellow]⚠[/yellow] Warnings:")
                for warning in result["warnings"]:
                    console.print(f"  - {warning}")
    except Exception as e:
        console.print(f"[red]Error:[/red] {e}")

if __name__ == "__main__":
    cli()
