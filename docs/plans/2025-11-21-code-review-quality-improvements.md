# Code Review Quality Improvements

**Date:** 2025-11-21
**Status:** Implemented

## Problem Statement

The Claude code review workflow was generating false positives and nitpicky feedback. Claude often tried too hard to find issues, leading to review comments that were either flat-out wrong or negligible, which eroded trust in the automated review process.

## Solution: Quality-First Review Philosophy

We implemented a comprehensive quality-first approach that prioritizes **accuracy over volume**. The core principle: no feedback is better than incorrect feedback.

## Key Changes

### 1. Philosophy Shift

The review prompts now explicitly:
- State that accuracy matters more than finding issues
- Give Claude permission to approve good code
- Emphasize that false positives harm trust more than missed issues
- Encourage highlighting what's done well when code is solid

### 2. Confidence-Gated Severity System

Issues are categorized with required confidence thresholds:

**CRITICAL (95%+ confidence required)** - Must fix before merge:
- Data races, goroutine leaks, or concurrency issues
- Nil pointer dereferences or panics
- Security vulnerabilities (injection, improper input validation)
- Logic errors causing incorrect behavior
- Breaking API changes without migration path

**IMPORTANT (90%+ confidence required)** - Should fix before merge:
- Inefficient algorithms with measurable performance impact
- Missing error handling or resource leaks (unclosed files, etc.)
- Improper use of defer, channels, or goroutines
- Architectural violations reducing maintainability
- Deviation from project patterns in CLAUDE.md

**MINOR (100% confidence required)** - Only if genuinely beneficial:
- Using more appropriate standard library functions (e.g., slices package)
- Simplifications improving readability without changing behavior
- Edge case handling that's currently missing but clearly should exist

**DON'T REPORT:**
- Style preferences or formatting (gofmt handles this)
- Micro-optimizations without measurable impact
- Anything uncertain or just "different from preference"

### 3. Go-Specific Review Criteria

The prompts now include focused guidance on Go best practices and architectural patterns:

- **Error handling**: Proper error wrapping, not swallowing errors, returning vs logging
- **Nil safety**: Potential nil pointer dereferences in pointers, maps, interfaces
- **Resource management**: Correct defer usage for cleanup
- **Concurrency**: Data races, channel usage, goroutine leaks
- **Standard library preference**: Use standard library over reinventing
- **Architecture**: Proper package boundaries (internal/stats, internal/reader, internal/display)
- **Type safety**: Appropriate use of Go's type system

### 4. Structured Output Format

Reviews now follow consistent templates:

**When issues are found:**
```
## Code Review

I found [X] issues requiring attention:

### CRITICAL
- `path/to/file.go:42` - [Specific issue with explanation and suggested fix]

### IMPORTANT
- `path/to/file.go:156` - [Specific issue with explanation and suggested fix]

### MINOR (optional improvements)
- `path/to/file.go:89` - [Clear benefit explained]

---
All issues verified against CLAUDE.md conventions and codebase context.
```

**When code looks good:**
```
## Code Review

This PR looks solid. No critical or important issues found.

**What's done well:**
- [Specific positive observation]
- [Another strength]

The implementation aligns with Go idioms and project architectural patterns.
```

### 5. Workflow-Specific Implementations

**claude-code-review.yml (Automated Reviews):**
- Runs automatically on all PRs
- Uses comprehensive quality-first prompt
- Focuses on critical issues, Go idioms, and architectural patterns

**claude.yml (Interactive Reviews):**
- Triggered by @claude mentions in comments
- Kept flexible for ad-hoc questions and tasks
- Added usage examples showing how to request quality-focused reviews
- Includes optional quality-first prompt (commented out) that users can enable

## Benefits

1. **Reduced false positives:** Claude only reports issues it's confident about
2. **Higher trust:** Fewer incorrect or negligible comments means reviews are more valuable
3. **Positive reinforcement:** Good code gets acknowledged, not nitpicked
4. **Clearer priorities:** Severity levels help prioritize what actually needs fixing
5. **Better Go alignment:** Specific focus on Go idioms and this project's patterns

## Implementation Files

- `.github/workflows/claude-code-review.yml` - Automated PR reviews with quality-first prompt
- `.github/workflows/claude.yml` - Interactive @claude mentions with usage guidance
- `docs/plans/2025-11-21-code-review-quality-improvements.md` - This design document

## Future Considerations

- Monitor review quality over the next few PRs
- Adjust confidence thresholds if needed based on actual results
- Consider adding examples of "good" vs "bad" issues in the prompt if patterns emerge
- Potentially expand Go-specific criteria as the codebase evolves
