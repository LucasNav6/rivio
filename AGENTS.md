# Coding Conventions

- Write all code comments, documentation comments, logs, and errors in English.
- Document every function with its responsibility, inputs, and return values.
- Document every type with its purpose and role in the application.
- Add ordered `// 1.`, `// 2.`, and subsequent comments for steps inside functions.
- Use early returns to handle errors and invalid input close to where they occur.
- Use structured logs that identify Rivio, the intended action, and its result.
- Never log API keys or other credentials.
- Keep provider-specific construction in `internal/provider`; keep Git operations in `internal/git`.
