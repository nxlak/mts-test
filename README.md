# Go Dependency Update Checker

## Installation

```bash
make build
```

## Usage

```bash
# Check a repository for dependency updates
./mts-test github.com/stretchr/testify
# or
./mts-test https://github.com/nxlak/test

# Show all dependencies (including indirect)
./mts-test -all github.com/example/repo

# Output as JSON
./mts-test -format json github.com/example/repo
```

## Flags

```
-format string
    Output format: text | json (default: text)

-all
    Report both direct and indirect dependencies (default: false)
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