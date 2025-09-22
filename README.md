# Ncobase CLI

Command-line tool for generating Ncobase projects and components.

## Installation

### From Source

```bash
go install github.com/ncobase/cli@latest
```

### From Release

Download the appropriate binary from the releases page and place it in your PATH.

## Usage

```bash
# Show help
nco --help

# Create a new core module
nco create core auth-service

# Create a business module with MongoDB
nco create business payment --use-mongo --with-test

# Create a plugin
nco create plugin logger

# Create a standalone application
nco create webapp dashboard --standalone

# Show version information
nco version
```

## Commands

- `create` - Generate new components (core, business, plugin, or custom)
- `version` - Show version information
- `docs` - Generate documentation
- `plugin` - Plugin management commands
- `migrate` - Database migration commands

## Development

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Install locally
make install

# Run tests
make test

# Clean build artifacts
make clean
```

## License

See [LICENSE](LICENSE) file for details.
