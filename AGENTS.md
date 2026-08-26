# Agents Guide for Prometheus OTLP Translator

Use this guide to align contributions with the conventions expected in this repository.

## Pull Requests

- Use `area: short description` titles, for example `LabelNamer: add Build method`.
- Link fixes with a GitHub closing keyword, for example `Fixes #123`.

## Commits

- Each commit must compile and pass tests independently, except when one commit adds a test that exposes a bug and the next fixes it.
- Keep commits small and focused. Do not bundle unrelated changes.
- Sign off every commit with `git commit -s` to satisfy the DCO requirement.
- Do not include unrelated local changes in the pull request.

## Testing

- Run Go tests without `-count=1` by default so Go can reuse cached results. Use `-count=1` only when uncached execution is specifically required, and state why.
- Add a regression test for every bug fix.
- Add unit or end-to-end coverage for new behavior and exported API changes.
- Use realistic data and behavior in tests.
- Use exported APIs where possible to keep tests close to real library usage.
- Prefer extending existing table-driven tests over adding new test functions. Adapt an existing test to a table-driven form when that avoids duplicated cases.

## Performance Work

- Add a benchmark that demonstrates every claimed performance improvement.
- Compare before and after results with `go test -run='^$' -count=6 -bench=. -benchmem ./...` and include `benchstat` output in the pull request.
- Address benchmark regressions or explain why the affected case is not important.
- Reuse allocations in hot paths where appropriate.
- When reusing buffers passed through interfaces, document that callers must copy the contents and must not retain references.
- Link supporting analysis for complex performance changes.

## Code and Documentation

- Follow [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments) and the formatting and style guidance in [Go: Best Practices for Production Environments](https://peter.bourgon.org/go-in-production/#formatting-and-style).
- State assumptions that affect correctness or design.
- Document important ownership and lifetime semantics at the interface definition, not only in the implementation.
- Add doc comments for all exported declarations.
- Start code comments with a capital letter and end them with a full stop.
- Prefer compact doc comments focused on the essential contract, invariants, and non-obvious behavior. Keep implementation details in local comments near the code they explain, and avoid duplicating them in declaration comments.
- Run `make lint` before submitting changes. Fix findings instead of suppressing them unless they are clear false positives.
- Use `//nolint:linter1[,linter2,...]` sparingly and name each suppressed linter.

## Scope Discipline

- Keep unrelated changes in separate pull requests.
- Put preparatory refactors in separate commits.
- Split large changes into preparatory and follow-up pull requests, linking them with `Part of #NNN` or `Depends on #NNN`.

## CI and Workflow Changes

- Declare required GitHub token permissions explicitly in workflow files. Missing permissions can otherwise surface as misleading authorization failures.
