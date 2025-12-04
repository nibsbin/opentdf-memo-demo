# Masterprompt - Agent Template Generator

Quick and dirty tool to generate user agent prompt files from a YAML configuration.

## Files

- `users.yaml` - User configuration (names, roles, credentials, access)
- `render_agents.py` - Python script to render agent markdown files

## Usage

```bash
# Render to default location (../.github/agents/)
python render_agents.py

# Render to a custom directory
python render_agents.py --output ./test-output

# Dry run - print to stdout without writing files
python render_agents.py --dry-run

# Use a different config file
python render_agents.py --config other_users.yaml
```

## Configuration

Edit `users.yaml` to add/modify users. Each user needs:

- `name` - Full name with rank (e.g., "Maj Evan Riley")
- `client_id` - OAuth client ID for Keycloak
- `role` - Position title (e.g., "KC-46 Aircraft Commander")
- `access_summary` - Markdown-formatted access description
- `access` - Structured access data (flights, clearance, maintenance)

Defaults (client_secret, tools, common_tasks, restrictions) are shared across all users but can be overridden per-user.

## Adding a New User

1. Add entry to `users.yaml` under `users:`
2. Run `python render_agents.py`
3. New agent file appears in `.github/agents/`
