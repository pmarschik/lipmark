# lipmark

## Build & Test Commands

- `mise run setup` -- install dependencies
- `mise run check` -- run all quality gates (format + lint + test)
- `mise run fmt` -- format code
- `mise run lint` -- run linters
- `mise run test` -- run tests
- `mise run build` -- build project

## Conventions

### Commits

Use Conventional Commits strictly:

    <type>(<scope>): <description>

Types: feat, fix, refactor, build, ci, chore, docs, style, perf, test
Scopes: defined in `cog.toml` -- update that file when adding new scopes.

Every commit MUST follow this format. The CI pipeline enforces this via git-cliff.

### API Stability

This is a public Go library. Breaking changes affect downstream consumers.

- **NEVER introduce breaking API changes without asking the user first**
- Breaking changes MUST use `feat!:` or `fix!:` commit prefix (triggers major version bump)
- Always try to maintain backward compatibility: add new functions/types instead of changing existing ones
- Deprecate before removing: mark old APIs with `// Deprecated:` and keep them for at least one minor version
- Adding new exported functions, types, or methods is NOT breaking
- Changing function signatures, removing exports, or changing behavior IS breaking

### Code Quality

- Run `mise run check` before pushing
- All linters must pass with zero warnings
- Tests must pass
- Keep README.md up to date when behavior or API changes

### Version Control

- Primary VCS: jj (jujutsu)
- Run `mise run check` before `jj git push`
- Do not push directly -- prompt the user (hardware key signing)
