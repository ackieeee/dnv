# dnv

`dnv` is a Go-based CLI that compares `KEY=VALUE` style files such as `.env` and prints differences in a diff-like layout. It helps you validate configuration drift across environments with minimal effort.

## Features

- **Diff-style output**: Uses `---/+++` headers and `+/- key=value` lines to highlight changes clearly.
- **Missing key detection**: Flags keys that appear in only one of the files.
- **Single static binary**: Ships as a standalone Go binary without additional runtime dependencies.

## Installation

```bash
go install github.com/ackieeee/dnv@latest
```

During local development you can also build the binary with `go build .` from the repository root.

## Usage

```bash
dnv <first-file> <second-file>
```

When no differences are detected, the command prints All keys and values match.. If differences exist, it prints the diff output.

### Example output

```bash
dnv sample1.txt sample2.txt

--- sample1.txt
+++ sample2.txt
- NAME=staging
- API_URL=https://stg.example.com
+ NAME=production
+ API_URL=https://api.example.com
```

## Development notes

Issues and pull requests are welcome for bug reports or improvement ideas.
