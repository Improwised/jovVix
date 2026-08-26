# Contributing to jovVix

Thank you for your interest in contributing to jovVix! This guide will help you get started.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [How to Contribute](#how-to-contribute)
- [Branch Naming Convention](#branch-naming-convention)
- [Commit Message Convention](#commit-message-convention)
- [Pull Request Process](#pull-request-process)
- [Running Tests](#running-tests)
- [Code Style](#code-style)
- [Reporting Bugs](#reporting-bugs)
- [Requesting Features](#requesting-features)
- [Contributing to Documentation](#contributing-to-documentation)
- [Semantic Versioning](#semantic-versioning)

## Code of Conduct

This project adheres to the [Contributor Covenant Code of Conduct](./CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## Getting Started

1. **Find an issue:** Look for issues labeled [`good-first-issue`](https://github.com/Improwised/jovVix/labels/good-first-issue) or [`help-wanted`](https://github.com/Improwised/jovVix/labels/help-wanted)
2. **Comment on the issue:** Let others know you're working on it
3. **Fork the repository:** Click the Fork button on GitHub
4. **Clone your fork:**
   ```bash
   git clone https://github.com/<your-username>/jovVix.git
   cd jovVix
   ```
5. **Create a branch:** See [Branch Naming Convention](#branch-naming-convention)
6. **Make your changes**
7. **Submit a pull request**

## How to Contribute

### Bug Fixes

1. Check [existing issues](https://github.com/Improwised/jovVix/issues) for the bug
2. If not reported, [create a bug report](https://github.com/Improwised/jovVix/issues/new/choose)
3. Create a branch from `develop` with `fix/` prefix
4. Fix the bug and add tests if applicable
5. Submit a pull request

### New Features

1. [Open a feature request](https://github.com/Improwised/jovVix/issues/new/choose) first to discuss the feature
2. Wait for maintainer approval before starting work
3. Create a branch from `develop` with `feat/` prefix
4. Implement the feature with tests
5. Submit a pull request

### Documentation

1. Check [documentation issues](https://github.com/Improwised/jovVix/labels/type/documentation)
2. Create a branch from `develop` with `docs/` prefix
3. Make your documentation changes
4. Submit a pull request

## Branch Naming Convention

Use the format: `<type>/<short-description>`

| Type | Use Case | Example |
|------|----------|---------|
| `feat/` | New feature | `feat/quiz-timer-display` |
| `fix/` | Bug fix | `fix/scoreboard-calculation` |
| `docs/` | Documentation | `docs/api-setup-guide` |
| `refactor/` | Code refactoring | `refactor/quiz-service-structure` |
| `test/` | Adding tests | `test/quiz-controller-unit` |
| `chore/` | Maintenance tasks | `chore/update-dependencies` |
| `ci/` | CI/CD changes | `ci/add-coverage-reporting` |

## Commit Message Convention

We use [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

**Examples:**
```
feat(quiz): add timer display for timed quizzes
fix(reports): correct score calculation for negative points
docs(readme): update installation instructions
test(api): add unit tests for quiz creation endpoint
chore(deps): update Go dependencies to latest versions
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Adding or updating tests
- `refactor`: Code refactoring (no feature change)
- `chore`: Maintenance tasks
- `ci`: CI/CD configuration changes
- `style`: Code style changes (formatting, no logic change)

## Pull Request Process

1. **Fill out the PR template** completely
2. **Link the issue:** Use `Closes #123` to link to the related issue
3. **Ensure all checks pass:** Tests, linting, and build must be green
4. **Request a review:** From a maintainer or collaborator
5. **Address feedback:** Make requested changes and push new commits
6. **Merge:** A collaborator will merge your PR once approved

### PR Checklist

- [ ] Code follows the project's style guidelines
- [ ] Tests pass locally (`npm run test` for frontend, `go test ./...` for backend)
- [ ] No new warnings or errors
- [ ] Documentation updated (if applicable)
- [ ] PR has a clear, descriptive title
- [ ] PR links to related issue(s)

## Running Tests

### Frontend Tests

```bash
cd app
npm run test           # Run all tests
npm run test -- --watch  # Watch mode
```

### Backend Tests

```bash
cd api
go test ./...          # Run all tests
go test -v ./...       # Verbose output
go test -cover ./...   # With coverage
```

### Linting

**Frontend:**
```bash
cd app
npm run lint           # Check for issues
npm run lint-fix       # Auto-fix issues
```

**Backend:**
```bash
cd api
golangci-lint run      # Run linter
```

## Code Style

- **Go:** Follow standard Go conventions. We use `golangci-lint` for enforcement.
- **Vue/TypeScript:** Follow the ESLint configuration in the project. Use `npm run lint-fix` to auto-format.
- **General:** Write clear, self-documenting code. Add comments for complex logic.

## Reporting Bugs

Use the [Bug Report template](https://github.com/Improwised/jovVix/issues/new/choose) and include:

- Steps to reproduce
- Expected behavior
- Actual behavior
- Environment details (OS, browser, versions)
- Screenshots if applicable

## Requesting Features

Open a [Feature Request](https://github.com/Improwised/jovVix/issues/new/choose) and include:

- Problem statement
- Proposed solution
- Alternatives considered

## Contributing to Documentation

Documentation improvements are always welcome! Check for:

- Typos and grammar errors
- Missing setup instructions
- Outdated information
- Missing examples

Open a [new issue](https://github.com/Improwised/jovVix/issues/new/choose) with the appropriate template to report doc issues.

## Semantic Versioning

This project follows [Semantic Versioning](https://semver.org/) (semver):

- **MAJOR** (`X.0.0`): Incompatible API changes (breaking changes to `/api/v1/` endpoints, database schema breaking changes, removal of features)
- **MINOR** (`0.X.0`): New functionality in a backwards-compatible manner (new endpoints, new question types, new UI features)
- **PATCH** (`0.0.X`): Backwards-compatible bug fixes (correct score calculations, UI fixes, documentation updates)

### Version Tags

Use the `v` prefix: `v1.2.3`, `v2.0.0-beta.1`

```bash
git tag -a v1.2.3 -m "Release v1.2.3"
git push origin v1.2.3
```

### Pre-release Versions

Use suffixes for pre-releases: `v1.2.3-beta.1`, `v1.2.3-rc.1`

### What Counts as a Breaking Change

- Removing or renaming API endpoints
- Changing request/response JSON structure of existing endpoints
- Changing database column types or removing columns
- Removing environment variables that were previously required
- Changing authentication/authorization behavior

---

Thank you for contributing to jovVix! 🎉
