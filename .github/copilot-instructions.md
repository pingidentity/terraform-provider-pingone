# Copilot Custom Instructions

## Goal

This project is a Terraform provider written in Go using the Terraform Plugin Framework. The goal is to implement high-quality, idiomatic, and maintainable provider logic, following official HashiCorp guidelines.

## Coding Style

- Follow idiomatic Go patterns.
- Keep functions small, readable, and testable.
- Use context-aware functions (`context.Context`) and return `diag.Diagnostics` for errors or warnings.
- Prefer composition over inheritance.
- Avoid global variables; use struct-based dependency injection if needed.
- Write helpful comments for exported functions and types.

## Libraries and Frameworks

- Use the Terraform Plugin Framework
- Use `schema.Schema`, `resource.Resource`, and `data_source.DataSource` constructs from the Plugin Framework.
- Use `tfsdk` package conventions as shown in HashiCorp documentation.

## Things to Avoid

- Do not use or suggest deprecated Terraform SDK v2 patterns.
- Avoid suggesting code that bypasses Terraform’s type system or lifecycle model.
- Avoid global/shared mutable state unless strictly necessary.

## Testing

- Use the `testing` package and table-driven tests.
- Prefer mocks or fakes over real Terraform infrastructure calls in unit tests.
- Follow Go testing conventions and maintain high test coverage.
- Tests should be comprehensive and follow best practices.

## Documentation

- Accurately comment all public structs, functions, and methods using Go doc style.
- Include usage examples where appropriate.

## Interaction Guidelines

- Provide clear, concise, and factual responses.
- Avoid unnecessary elaboration, apologizing or making conciliatory statements.
- Avoid conversational statements and hyperbole
- Code snippets should be complete and functional examples
- When explaining a code snippet, focus on the key aspects relevant to the task.
- When refactoring code, maintain the original functionality while improving readability or performance.
- Documentation should be clear, concise, and follow Go doc conventions.
- Generated code must adhere to the project's coding style and conventions.
