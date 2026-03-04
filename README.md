# scaffold

A meta scaffolding CLI for AI agent contexts. Define agents and skills once,
then activate them to automatically populate context files like `CLAUDE.md` or
`.cursorrules` that your AI tools pick up automatically.

## How It Works

1. Define agents in `agents.md` — each agent has a name, a list of skills, and a system prompt.
2. Write skill files in `AgentSkills/` — reusable instruction blocks (e.g. `coding.md`, `testing.md`).
3. Run `scaffold use <agent-name>` — scaffold combines the agent prompt with its skill files and writes the result to your configured target file(s).

## Installation

Requires [Go 1.22+](https://go.dev/dl/).

```bash
git clone https://github.com/msallean-git/scaffold.git
cd scaffold
go build -o scaffold .
```

Or on Windows:

```bat
build.bat
```

## Quick Start

```bash
# Initialize scaffold in your project
scaffold init

# List available agents
scaffold list

# Activate an agent (writes to CLAUDE.md by default)
scaffold use general-dev

# Preview output without writing any files
scaffold use general-dev --dry-run

# Clear the active agent context
scaffold reset
```

## Commands

### `scaffold init`

Sets up scaffold in the current directory. Creates:
- `scaffold.config.json` — configuration file
- `agents.md` — agent definitions (includes a sample agent)
- `AgentSkills/` — directory for skill files
- `AgentSkills/coding.md` — a starter skill file

```bash
scaffold init                        # default target: CLAUDE.md
scaffold init --target .cursorrules  # use a different target file
scaffold init --force                # overwrite existing files
```

### `scaffold list`

Displays all agents defined in `agents.md`.

```bash
scaffold list          # table output
scaffold list --json   # JSON output for scripting
```

### `scaffold create <agent-name>`

Adds a new agent block to `agents.md`.

```bash
scaffold create backend-dev --skills coding,testing,database
scaffold create frontend-dev --skills coding,react --instructions "You are a frontend developer..."
scaffold create devops --skills docker,ci --create-skills  # creates stub files for missing skills
```

### `scaffold use <agent-name>`

Activates an agent: loads its skills, renders the combined output, and writes
it to all configured target files.

```bash
scaffold use backend-dev
scaffold use backend-dev --dry-run   # preview without writing
scaffold use backend-dev --verbose   # show extra detail
```

### `scaffold reset`

Clears the active agent context from all target files and unsets `activeAgent`
in the config.

```bash
scaffold reset
```

## Configuration

`scaffold.config.json` is created by `scaffold init` and controls all behaviour:

```json
{
  "version": 1,
  "targets": [
    { "path": "CLAUDE.md", "mode": "overwrite" }
  ],
  "agentsFile": "agents.md",
  "skillsDir": "AgentSkills",
  "activeAgent": ""
}
```

### Target modes

| Mode | Behaviour |
|------|-----------|
| `overwrite` | Replaces the entire target file with the rendered output. |
| `section` | Replaces only the content between `<!-- scaffold:start -->` and `<!-- scaffold:end -->` markers, leaving the rest of the file untouched. |

You can define multiple targets to write to several files at once:

```json
"targets": [
  { "path": "CLAUDE.md",      "mode": "overwrite" },
  { "path": ".cursorrules",   "mode": "section"   }
]
```

## agents.md Format

```markdown
# Agents

## backend-dev

**Skills:** coding, testing, database

You are a senior backend developer. Write clean, well-tested code...

## frontend-dev

**Skills:** coding, react, accessibility

You are a frontend developer focused on user experience...
```

- `## <name>` — starts an agent definition
- `**Skills:** skill1, skill2` — comma-separated list of skill file basenames (without `.md`)
- Everything else until the next `##` or end of file is the agent's system prompt

## Output Format

The rendered output written to target files looks like this:

```markdown
<!-- scaffold:start agent=backend-dev -->
# Agent: backend-dev

## Instructions

You are a senior backend developer...

## Skills

### coding

(contents of AgentSkills/coding.md)

### testing

(contents of AgentSkills/testing.md)

<!-- scaffold:end -->
```

## License

MIT
