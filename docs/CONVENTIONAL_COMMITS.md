# Conventional commits

A guide to writing consistent, semantic commit messages that enable automatic changelog generation.

## Format

```text
<type>(<scope>): <subject>

<body>

<footer>
```

## Types

- **feat**: New feature
- **fix**: Bug fix
- **docs**: Documentation changes
- **style**: Code style changes (formatting, missing semi-colons, etc.)
- **refactor**: Code refactoring (neither fixes a bug nor adds a feature)
- **perf**: Performance improvements
- **test**: Adding or updating tests
- **chore**: Maintenance tasks (updating dependencies, etc.)
- **ci**: CI/CD changes
- **build**: Build system changes

## Examples

### Simple Feature

```text
feat: Add user authentication
```

### Feature with Scope

```text
feat(api): Add user registration endpoint
```

### Bug Fix

```text
fix: Correct off-by-one error in pagination
```

### Breaking Change

```text
feat!: Redesign API response format

BREAKING CHANGE: The API now returns data in a nested structure.
Migration guide: https://example.com/migration
```

### With Body and Footer

```text
feat(auth): Add OAuth2 support

Implement OAuth2 authentication flow with support for
Google and GitHub providers.

Closes #123
```

## Scope

The scope is optional and provides additional context:

- **api**: API changes
- **cli**: CLI changes
- **config**: Configuration changes
- **deps**: Dependency updates
- **ui**: User interface changes

## Subject

- Use imperative, present tense: "add" not "added" or "adds"
- Don't capitalize the first letter
- No period (.) at the end
- Keep it under 50 characters

## Body

- Optional, provides detailed explanation
- Wrap at 72 characters
- Explain what and why, not how
- Separate from subject with a blank line

## Footer

- Optional, references issues or breaking changes
- `Closes #123` - Closes an issue
- `Refs #456` - References an issue
- `BREAKING CHANGE:` - Describes breaking changes

## Why Conventional Commits?

1. **Automatic Changelog** - Tools can generate changelogs from commit messages
2. **Semantic Versioning** - Automatically determine version bumps
3. **Better History** - Clear, searchable commit history
4. **Team Communication** - Consistent format helps team understanding

## GoReleaser Integration

GoReleaser can automatically group commits by type in release notes:

```yaml
changelog:
  groups:
    - title: Features
      regexp: '^.*?feat(\([[:word:]]+\))??!?:.+$'
    - title: Bug Fixes
      regexp: '^.*?fix(\([[:word:]]+\))??!?:.+$'
```

## Tools

### Commitizen

Interactive tool for creating conventional commits:

```bash
npm install -g commitizen cz-conventional-changelog
git cz
```

### commitlint

Validates commit messages:

```bash
npm install -g @commitlint/cli @commitlint/config-conventional
```

## Resources

- [Conventional Commits Specification](https://www.conventionalcommits.org/)
- [Angular Commit Guidelines](https://github.com/angular/angular/blob/main/CONTRIBUTING.md#commit)
- [Commitizen](https://github.com/commitizen/cz-cli)
