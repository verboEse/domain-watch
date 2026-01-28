# Project Context

## Purpose
Domain Watch is a CLI application that monitors domain registration changes and expiration dates. It fetches public WHOIS records on a configurable schedule and sends notifications when:
- Domain expiration date passes a configured threshold
- Domain status code changes

The tool supports multiple notification providers (Gotify, Telegram) and includes Prometheus metrics for monitoring.

## Tech Stack
- **Language**: Go 1.25.6
- **CLI Framework**: Cobra (spf13/cobra)
- **Notifications**: Telegram (go-telegram/bot), Gotify (gotify/server)
- **WHOIS Lookup**: likexian/whois and likexian/whois-parser
- **Monitoring**: Prometheus client (prometheus/client_golang)
- **Logging**: slog with tint formatter (lmittmann/tint)
- **Utilities**: Custom utilities (gabe565.com/utils)
- **Testing**: testify assertions
- **Containerization**: Docker, Docker Compose
- **Distribution**: goreleaser for releases
- **Deployment**: systemd service file

## Project Conventions

### Code Style
- **Language**: Go with standard Go conventions (gofmt)
- **Package Structure**: Internal packages under `internal/` for private code
- **Logging**: Use `slog` package with structured logging
- **Error Handling**: Return errors from functions, log errors in main flow
- **Naming**: CamelCase for Go functions and variables, UPPERCASE for constants
- **Comments**: Document exported functions and packages

### Architecture Patterns
- **Modular Design**: Organized into functional packages:
  - `config/`: Configuration loading, CLI flags, completions
  - `domain/`: Domain data structures and WHOIS operations
  - `integration/`: Notification provider integration (Gotify, Telegram)
  - `message/`: Notification message formatting
  - `metrics/`: Prometheus metrics collection
  - `util/`: Helper utilities and error handling
- **Context-based Configuration**: Uses Go context to pass configuration through the application
- **Dependency Injection**: Integrations are initialized based on config
- **Cobra Commands**: CLI command structure with flag registration and validation

### Testing Strategy
- Unit tests alongside source files with `_test.go` suffix
- Integration tests for notification providers (setup files included)
- Use `stretchr/testify` for assertions
- Test coverage for domain parsing and message generation

### Git Workflow
- Main branch is the default
- Follow standard GitHub workflow
- Releases managed through goreleaser

## Domain Context
**WHOIS Protocol Knowledge**:
- Domain expiration dates are fetched from public WHOIS records
- Status codes (e.g., active, expired, pending delete) can change
- Different registrars may have different WHOIS response formats

**Notification Context**:
- Threshold configuration: Default is [1, 7] days before expiration
- Multiple domains can be monitored in a single run
- Sleep duration between domain checks prevents rate limiting

**Metrics Context**:
- Prometheus metrics endpoint for monitoring check success/failure rates
- Helps track service health and domain status changes

## Important Constraints
- Public WHOIS API rate limiting (3-second default sleep between checks)
- WHOIS record parsing variability across different registrars
- Notification provider credentials must be securely configured
- Docker deployment requires appropriate permission handling

## External Dependencies
- **WHOIS Servers**: Public WHOIS infrastructure for domain queries
- **Telegram Bot API**: Requires valid bot token and chat ID
- **Gotify Server**: Self-hosted notification server (URL and token required)
- **Prometheus**: For metrics collection and monitoring
