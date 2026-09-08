# Governance

This document describes how the jovVix project is governed, how decisions are made, and how community members can contribute to the project's direction.

## Roles

### Contributors

Anyone who submits a pull request, reports an issue, or participates in discussions. Contributors are encouraged to:

- Follow the [Contributing Guidelines](./CONTRIBUTING.md)
- Adhere to the [Code of Conduct](./CODE_OF_CONDUCT.md)
- Be respectful and constructive in all interactions

### Collaborators

Contributors who have demonstrated sustained commitment and expertise. Collaborators are invited by existing maintainers based on:

- Consistent, high-quality contributions
- Deep understanding of the codebase
- Active participation in code reviews and issue triage
- Alignment with the project's values and vision

Collaborators receive:

- Commit access to the repository
- Ability to merge pull requests
- Input on project direction

### Maintainers

The core team responsible for the project's long-term direction. Maintainers have final say on:

- Architectural decisions
- Release planning and versioning
- Security vulnerability responses
- Governance changes

## Decision Making

### Day-to-Day Decisions

- Code style, minor refactors, bug fixes: Collaborator approval (1 reviewer)
- New features, significant changes: 2 collaborator approvals required
- Breaking changes, architectural shifts: Maintainer approval required

### Consensus Process

1. **Proposal:** Open a GitHub Issue or Discussion with a clear proposal
2. **Discussion:** Community feedback period (minimum 7 days for major changes)
3. **Decision:** Maintainers make final decision, documented in the issue
4. **Implementation:** Approved changes are implemented via pull request

### Voting

For major decisions that affect the entire community:

- Maintainers have 3 votes each
- Collaborators have 1 vote each
- Simple majority wins
- Voting period: 14 days minimum
- Abstentions do not count toward the total

## Adding a Collaborator

1. An existing maintainer nominates a contributor
2. The nomination is discussed by maintainers (private channel)
3. If approved, the contributor is invited via GitHub
4. The new collaborator is announced in the project's communication channels

## Removing a Collaborator

Collaborators may be removed for:

- Sustained inactivity (6+ months without contributions or reviews)
- Violation of the Code of Conduct
- Actions harmful to the project

The process:
1. A maintainer raises the concern privately with the collaborator
2. If not resolved, a maintainer vote is held
3. A 2/3 majority is required for removal
4. The collaborator is notified and given an opportunity to respond

## Changes to Governance

Changes to this governance document require:

- Maintainer approval
- 14-day public comment period
- 2/3 supermajority vote of maintainers

## Communication

- **GitHub Issues & PRs:** Primary channel for technical discussions
- **GitHub Discussions:** For general questions and proposals
- **Security reports:** Via email (see [SECURITY.md](./SECURITY.md))
