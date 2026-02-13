# Agent Workflow Guidelines

## Overview

This document defines the standard workflow for AI agents working on the OpenPDF Reader project.

---

## Task Management

### TODO.md Structure

All work items are tracked in `TODO.md` at the project root.

Format:
```markdown
# TODO

## In Progress
- [ ] Task description

## Completed
- [x] Completed task description (commit: abc1234)
```

### Workflow

1. **Before starting work:**
   - Read `TODO.md` to understand current tasks
   - Identify the next task to work on
   - Mark task as in-progress if using internal todo tracking

2. **During work:**
   - Make focused, incremental changes
   - Test changes when applicable
   - Keep commits atomic (one logical change per commit)

3. **After completing a task:**
   - run tests to confirm it's workingS
   - Update `TODO.md` - mark task as done with `[x]`
   - Move completed task to "Completed" section
   - Commit the change with a descriptive message

---

## Git Workflow

### Commit Guidelines

- **Format:** `<type>: <short description>`
- **Types:**
  - `feat` - New feature
  - `fix` - Bug fix
  - `docs` - Documentation changes
  - `refactor` - Code refactoring
  - `test` - Adding/updating tests
  - `chore` - Build, config, dependency updates

- **Examples:**
  ```
  feat: add PDF page thumbnail sidebar
  fix: resolve memory leak in page renderer
  docs: update installation instructions
  ```

### Commit Sequence

```bash
# Stage changes
git add <files>

# Commit with message
git commit -m "<type>: <description>"

# Update TODO.md and commit
git add TODO.md
git commit -m "chore: mark <task> as complete"
```

Or combine if the TODO update is part of the feature:
```bash
git add . 
git commit -m "feat: implement feature X"
```

---

## Code Standards

### Go Conventions

- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Write tests for new functionality
- Handle errors explicitly - no silent failures
- Document exported functions and types

### File Organization

- Keep related code together
- One package per directory
- Test files alongside source: `foo.go` → `foo_test.go`

---

## Documentation

- Update `README.md` for user-facing changes
- Update `project.md` for architectural decisions
- Add inline comments for complex logic
- Keep `CHANGELOG.md` updated for releases

---

## Checklist Before Marking Done

- [ ] Code compiles without errors
- [ ] Tests pass (if applicable)
- [ ] Changes committed with proper message
- [ ] `TODO.md` updated
- [ ] No debug code or temporary files left behind
