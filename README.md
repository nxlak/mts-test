# Go Dependency Update Checker

## Installation

```bash
make build
```

## Usage

```bash
# Check a repository for dependency updates
./mts-test github.com/example/repo
# or
./mts-test https://github.com/example/repo

# Show all dependencies (including indirect)
./mts-test github.com/example/repo -a

# Output as JSON
./mts-test github.com/example/repo -f json
```

## Flags

```
-f, --format string
    Output format: text | json (default: text)

-a, --all
    Report both direct and indirect dependencies (default: false)

-h, --help
    Show help message
```

## Testing

Run all tests:

```bash
make test
```

## Requirements

- Go 1.25 or higher
- Git installed

## How it works

1. Clones the repository to a temporary directory
2. Parses the `go.mod` file
3. Queries `proxy.golang.org` for latest versions
4. Shows available updates
5. Cleans up temporary files