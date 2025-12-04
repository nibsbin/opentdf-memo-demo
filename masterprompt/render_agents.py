#!/usr/bin/env python3
"""
Quick and dirty script to template and render user agent files.

Usage:
    python render_agents.py                    # Renders to ../.github/agents/
    python render_agents.py --output ./out     # Renders to custom directory
    python render_agents.py --dry-run          # Print to stdout, don't write files
"""

import argparse
from pathlib import Path
from string import Template
import yaml


# Default template path
DEFAULT_TEMPLATE_PATH = Path(__file__).parent / "agent_template.md"


def load_template(template_path: Path) -> str:
    """Load the template file."""
    with open(template_path, "r") as f:
        return f.read()


def load_config(config_path: Path) -> dict:
    """Load the YAML configuration file."""
    with open(config_path, "r") as f:
        return yaml.safe_load(f)


def render_user(user: dict, defaults: dict, template_content: str) -> str:
    """Render a single user's agent file content."""
    # Merge defaults with user-specific values
    client_secret = user.get("client_secret", defaults["client_secret"])
    common_tasks = user.get("common_tasks", defaults["common_tasks"])
    restrictions = user.get("restrictions", defaults["restrictions"])

    # Format common tasks as markdown list
    tasks_str = "\n".join(f"- {task}" for task in common_tasks)

    # Format restrictions
    restrictions_str = "\n".join(restrictions)

    # Clean up access_summary (remove trailing whitespace per line)
    access_summary = user["access_summary"].rstrip()

    # Clean up bio (remove trailing whitespace per line)
    bio = user.get("bio", "").rstrip()

    # Build template context
    context = {
        "name": user["name"],
        "role": user["role"],
        "client_id": user["client_id"],
        "client_secret": client_secret,
        "access_summary": access_summary,
        "bio": bio,
        "restrictions": restrictions_str,
        "common_tasks": tasks_str,
    }

    # Use Template for safe substitution
    template = Template(template_content)
    return template.substitute(context)


def get_output_filename(client_id: str) -> str:
    """Generate output filename from client_id."""
    return f"{client_id}.agent.md"


def main():
    parser = argparse.ArgumentParser(
        description="Render user agent files from YAML configuration"
    )
    parser.add_argument(
        "--config",
        type=Path,
        default=Path(__file__).parent / "users.yaml",
        help="Path to YAML config file (default: users.yaml)",
    )
    parser.add_argument(
        "--output",
        type=Path,
        default=Path(__file__).parent.parent / ".github" / "agents",
        help="Output directory for agent files (default: ../.github/agents/)",
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Print rendered content to stdout instead of writing files",
    )
    parser.add_argument(
        "--template",
        type=Path,
        default=DEFAULT_TEMPLATE_PATH,
        help="Path to template file (default: agent_template.md)",
    )
    args = parser.parse_args()

    # Load configuration and template
    config = load_config(args.config)
    defaults = config["defaults"]
    users = config["users"]
    template_content = load_template(args.template)

    print(f"Loaded {len(users)} users from {args.config}")
    print(f"Using template: {args.template}")

    # Ensure output directory exists
    if not args.dry_run:
        args.output.mkdir(parents=True, exist_ok=True)

    # Render each user
    for user in users:
        content = render_user(user, defaults, template_content)
        filename = get_output_filename(user["client_id"])

        if args.dry_run:
            print(f"\n{'='*60}")
            print(f"# {filename}")
            print("=" * 60)
            print(content)
        else:
            output_path = args.output / filename
            with open(output_path, "w") as f:
                f.write(content)
            print(f"  Wrote: {output_path}")

    if not args.dry_run:
        print(f"\nRendered {len(users)} agent files to {args.output}")


if __name__ == "__main__":
    main()
