# Project Guidance

## Project Structure

- This project is a Go command-line application. Follow standard Go formatting, naming, package, and error-handling conventions.
- Keep Cobra command wiring and user-facing CLI behavior in `cmd`.
- Keep YAML parsing and configuration types in `internal/config`.
- Keep Git inspection and diff operations in `internal/git`; do not modify the repository as part of review.
- Keep the shared review prompt in `internal/prompt`.
- Keep provider selection and provider-specific construction in `internal/provider`. Keep API and CLI integrations in their respective subpackages, and shared CLI process execution in `internal/provider/cli/runner`.

## Documentation and Comments

- Document exported packages, types, constants, variables, and functions in Go doc style. Explain purpose and relevant behavior; for functions, describe inputs and results when they are not clear from the signature.
- Add comments to unexported declarations or code blocks when they explain non-obvious behavior, constraints, or rationale. Do not require comments that merely restate the code.
- Use ordered `// 1.`, `// 2.`, and subsequent step comments only for multi-step logic where they make the flow easier to follow. Do not add them to simple functions or wrappers.
- Write comments and documentation in English.

## Errors, Logging, and Credentials

- Return errors to the caller and add context with `%w` when wrapping an underlying error. Validate required inputs near their use and return early on invalid input.
- Use the existing `charm.land/log/v2` logger for operational events. When logging, use English messages and structured fields that identify the action and its outcome; avoid logging the same failure at multiple layers.
- Never log API keys, tokens, or other credentials. Do not include secrets in error messages.

## Go Changes

- Keep dependencies and implementation patterns consistent with the existing Go modules and standard library where practical.
- Keep provider transport details inside provider packages and keep command orchestration in `cmd`.
- Format Go files with `gofmt` after editing them.
