# Semantic Versioning Guide for Gaia MCP Go

This document provides a complete guide to understanding and using semantic versioning in the Gaia MCP Go project.

## Table of Contents

- [What is Semantic Versioning?](#what-is-semantic-versioning)
- [Version Format](#version-format)
- [Version Management](#version-management)
- [Version Bumping](#version-bumping)
- [Release Workflows](#release-workflows)
- [Best Practices](#best-practices)
- [Examples](#examples)
- [Troubleshooting](#troubleshooting)

## What is Semantic Versioning?

Semantic Versioning (SemVer) is a versioning scheme that gives meaning to version numbers. It helps developers and users understand what changes are included in each release.

### Key Benefits

- **Predictability**: Users know what to expect from version updates
- **Dependency Management**: Tools can automatically determine compatibility
- **Clear Communication**: Version numbers convey the nature of changes
- **Professional Standard**: Widely adopted across the software industry

## Version Format

Semantic versions follow the `MAJOR.MINOR.PATCH` format:

```
1.2.3
│ │ │
│ │ └─ PATCH: Bug fixes (backward compatible)
│ └─── MINOR: New features (backward compatible)
└───── MAJOR: Breaking changes (not backward compatible)
```

### Additional Labels

- **Pre-release**: `1.0.0-alpha.1`, `1.0.0-beta.2`, `1.0.0-rc.1`
- **Build metadata**: `1.0.0+20231215.abc123`
- **Combined**: `1.0.0-beta.1+20231215.abc123`

### When to Increment Each Number

| Change Type                        | Version Bump | Example           |
| ---------------------------------- | ------------ | ----------------- |
| Bug fixes, security patches        | PATCH        | `1.0.0` → `1.0.1` |
| New features, deprecations         | MINOR        | `1.0.0` → `1.1.0` |
| Breaking changes, removed features | MAJOR        | `1.0.0` → `2.0.0` |

## Version Management

### Checking Current Version

```bash
# Check current version
make version

# Using the built binary
./bin/gaia-mcp-go --version
./bin/gaia-mcp-go -v
```

### Building with Version Information

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Install with version info
make install
```

## Version Bumping

### Automated Version Bumping (Recommended)

Use the provided Make commands for automated version management:

```bash
# Patch version bump (bug fixes)
make bump-patch     # 1.0.0 → 1.0.1

# Minor version bump (new features)
make bump-minor     # 1.0.0 → 1.1.0

# Major version bump (breaking changes)
make bump-major     # 1.0.0 → 2.0.0

# Dry run (see what would happen)
make bump-dry
```

### Manual Version Tagging

For more control, create tags manually:

```bash
# Create a regular release
git tag v1.0.0
git push origin v1.0.0

# Create a pre-release
git tag v1.0.0-alpha.1
git push origin v1.0.0-alpha.1

# Create with annotation
git tag -a v1.0.0 -m "Release v1.0.0 - Initial stable release"
git push origin v1.0.0
```

### Pre-release Versioning

Create pre-release versions for testing:

```bash
# Alpha versions (early testing)
git tag v1.0.0-alpha.1 && git push origin v1.0.0-alpha.1

# Beta versions (feature complete)
git tag v1.0.0-beta.1 && git push origin v1.0.0-beta.1

# Release candidates (final testing)
git tag v1.0.0-rc.1 && git push origin v1.0.0-rc.1
```

## Release Workflows

### 1. Feature Release Workflow

```bash
# Ensure you're on main and up to date
git checkout main
git pull origin main

# For new features
make bump-minor

# Verify the release
./bin/gaia-mcp-go --version
```

### 2. Hotfix Release Workflow

```bash
# Fix the critical issue
git checkout main
# ... make your bug fixes ...
git add .
git commit -m "fix: resolve critical security issue"

# Create hotfix release
make bump-patch
```

### 3. Pre-release Testing Workflow

```bash
# Create alpha for initial testing
git tag v2.0.0-alpha.1
git push origin v2.0.0-alpha.1

# After testing, create beta
git tag v2.0.0-beta.1
git push origin v2.0.0-beta.1

# After more testing, create release candidate
git tag v2.0.0-rc.1
git push origin v2.0.0-rc.1

# Finally, create the stable release
make bump-major
```

### 4. Development Workflow

```bash
# During development (no releases)
git checkout -b feature/new-feature
# ... develop ...
git commit -m "feat: add new feature"

# When ready to release
git checkout main
git merge feature/new-feature
make bump-minor  # or appropriate bump type
```

## Best Practices

### Version Planning Strategy

- **Start with `v0.1.0`** for initial development
- **Use `v1.0.0`** when your API is stable and ready for production
- **Pre-releases** for testing: `alpha` → `beta` → `rc` → stable
- **Patch releases** for hotfixes and bug fixes only

### Release Timing

- **Major releases**: Plan carefully, announce breaking changes
- **Minor releases**: Regular feature releases (monthly/quarterly)
- **Patch releases**: As needed for critical fixes
- **Pre-releases**: For testing new features before stable release

### Commit Message Conventions

Use conventional commit format to indicate the type of change:

```bash
feat: add new API endpoint          # → MINOR bump
fix: resolve authentication bug     # → PATCH bump
feat!: change API response format   # → MAJOR bump
docs: update documentation          # → No version bump
test: add unit tests               # → No version bump
```

### Documentation and Release Notes

Always update when releasing:

- **CHANGELOG.md** with detailed release notes
- **README.md** if there are API or usage changes
- **API documentation** for new features
- **Migration guides** for major version changes

## Examples

### Complete Feature Release

```bash
# 1. Ensure clean working directory
git status

# 2. Update documentation
vim CHANGELOG.md
vim README.md

# 3. Commit documentation changes
git add .
git commit -m "docs: update changelog and readme for v1.2.0"

# 4. Create the release
make bump-minor

# 5. Verify the release
./bin/gaia-mcp-go --version
```

### Emergency Hotfix Release

```bash
# 1. Create hotfix branch
git checkout -b hotfix/security-fix

# 2. Fix the critical issue
# ... make changes ...
git commit -m "fix: resolve critical security vulnerability"

# 3. Merge to main
git checkout main
git merge hotfix/security-fix

# 4. Create immediate patch release
make bump-patch

# 5. Clean up
git branch -d hotfix/security-fix
```

### Major Version Release with Pre-releases

```bash
# 1. Create alpha for early feedback
git tag v2.0.0-alpha.1
git push origin v2.0.0-alpha.1

# 2. Collect feedback, make changes, create alpha.2
git tag v2.0.0-alpha.2
git push origin v2.0.0-alpha.2

# 3. Feature freeze, create beta
git tag v2.0.0-beta.1
git push origin v2.0.0-beta.1

# 4. Final testing, create release candidate
git tag v2.0.0-rc.1
git push origin v2.0.0-rc.1

# 5. Final release
make bump-major
```

## Troubleshooting

### Common Version Management Issues

**❌ "Invalid semantic version format"**

```bash
# Check your current tags
git tag -l

# Tags should be: v1.2.3 (not 1.2.3 or version-1.2.3)
git describe --tags --abbrev=0
```

**❌ "Working directory not clean"**

```bash
# Check what files are modified
git status

# Commit or stash changes before version bumping
git add .
git commit -m "your message"
# or
git stash
```

**❌ Version bump didn't work**

```bash
# Check if you have the required scripts
ls -la scripts/bump-version.sh

# Ensure script is executable
chmod +x scripts/bump-version.sh

# Check Makefile has the bump targets
grep "bump-" Makefile
```

### Debugging Version Detection

```bash
# See what version would be detected
git describe --tags --always --dirty

# Check build-time version injection
make version

# Test version in your binary
go build && ./bin/gaia-mcp-go --version
```

### Recovery from Version Mistakes

**Wrong tag was created:**

```bash
# Delete the wrong tag locally
git tag -d v1.0.0

# Delete the wrong tag from remote
git push origin :refs/tags/v1.0.0

# Create the correct tag
git tag v1.0.1
git push origin v1.0.1
```

**Need to fix a released version:**

```bash
# Releases are immutable - create a new patch version
make bump-patch

# Or create a hotfix with specific fixes
git commit -m "fix: address issue in v1.0.0"
make bump-patch
```

**Accidentally bumped wrong version type:**

```bash
# If you haven't pushed yet, delete and recreate
git tag -d v2.0.0  # If you meant to do minor bump
git tag v1.1.0
git push origin v1.1.0

# If already pushed, just continue with the higher version
# and plan the next release accordingly
```

### Release Checklist

Before creating any release:

- [ ] **Tests pass**: `make test`
- [ ] **Documentation updated**: CHANGELOG.md, README.md
- [ ] **Clean working directory**: `git status`
- [ ] **Correct branch**: Usually `main` for releases
- [ ] **Version type correct**: patch/minor/major
- [ ] **Build works**: `make build`
- [ ] **Version displays correctly**: `./bin/gaia-mcp-go --version`

## Summary

This semantic versioning system helps you:

- ✅ **Communicate changes clearly** through version numbers
- ✅ **Automate release processes** with simple commands
- ✅ **Manage dependencies effectively** with semantic meaning
- ✅ **Plan releases strategically** with pre-release versions
- ✅ **Maintain professional standards** in software distribution

Remember: Good versioning is about **communication**. Each version number tells users what to expect, making your software more reliable and trustworthy.
