# Change: Fix golangci-lint configuration version

## Why
The golangci-lint CI job is failing with the error: `unsupported version of the configuration: ""`. This occurs because the `.golangci.yaml` file is missing the required `version` field that golangci-lint 2.8.0 expects. According to golangci-lint migration guide, all configurations must explicitly specify the configuration version to ensure compatibility.

## What Changes
- Add `version: 2` to the top of `.golangci.yaml` to specify the configuration version
- This enables golangci-lint 2.8.0 and newer to properly parse and validate the configuration
- The change maintains backward compatibility with existing linter settings

## Impact
- **Affected files**: `.golangci.yaml`
- **CI/CD**: Fixes the failing golangci-lint check in GitHub Actions workflows
- **Testing**: No code changes required; configuration-only update
- **Breaking changes**: None - existing linter configuration remains unchanged
