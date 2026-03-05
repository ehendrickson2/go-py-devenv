```markdown
# go-py-devenv

A small Golang utility that takes a repository and root directory location and sets up a Python environment using your choice of package manager.

## Prerequisites

- Go 1.16 or higher
- Git
- [UV](https://docs.astral.sh/uv/) package manager and/or [Poetry](https://python-poetry.org/)

## Setup

1. Clone this repository:
```bash
git clone https://github.com/yourusername/go-py-devenv.git
cd go-py-devenv
```

2. Build the project:
```bash
go build -o go-py-devenv
```

## Usage

Run the utility with an optional package manager flag:

```bash
./go-py-devenv [-pm uv|poetry]
```

### Flags

- `-pm` - Specifies the package manager to use. Options are `uv` (default) or `poetry`

The utility will prompt you for:
1. **Root directory** - The location where you want to clone the repository
2. **Repository URL** - The Git repository URL to clone (e.g., `https://github.com/user/repo.git`)

### UV (Default)

The utility will then:
- Clone the repository (or skip if it already exists)
- Create a UV virtual environment
- Install dependencies using either `uv sync` (for `pyproject.toml`) or `uv pip install` (for `requirements.txt`)

### Poetry

The utility will then:
- Clone the repository (or skip if it already exists)
- Install dependencies using `poetry install` (requires `pyproject.toml`)

## Examples

### Using UV (default):
```bash
./go-py-devenv
Root directory where you want the repository to be cloned to:
/home/user/projects
Repository you are cloning:
https://github.com/nautobot/nautobot.git
```

### Using Poetry:
```bash
./go-py-devenv -pm poetry
Root directory where you want the repository to be cloned to:
/home/user/projects
Repository you are cloning:
https://github.com/nautobot/nautobot.git
```
