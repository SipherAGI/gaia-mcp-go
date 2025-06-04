# Semantic Versioning Quick Reference

Essential commands and workflows for version management in Gaia MCP Go.

## 🚀 Quick Start

```bash
# Check current version
make version

# Bump version (choose one)
make bump-patch    # Bug fixes: 1.0.0 → 1.0.1
make bump-minor    # New features: 1.0.0 → 1.1.0
make bump-major    # Breaking changes: 1.0.0 → 2.0.0

# Preview before bumping
make bump-dry
```

## 📋 Version Management Commands

### Check Version

```bash
make version                          # Current version info
./bin/gaia-mcp-go --version          # Binary version
git describe --tags --always         # Git version
```

### Build with Version

```bash
make build                           # Build current platform
make build-all                       # Build all platforms
```

### Automated Bumping

```bash
make bump-patch                      # 1.0.0 → 1.0.1 (bug fixes)
make bump-minor                      # 1.0.0 → 1.1.0 (new features)
make bump-major                      # 1.0.0 → 2.0.0 (breaking changes)
make bump-dry                        # Preview without changes
```

### Manual Tagging

```bash
git tag v1.0.0                       # Create release tag
git push origin v1.0.0               # Push tag
git tag -a v1.0.0 -m "Release v1.0.0" # Annotated tag
```

## 🔄 Release Workflows

### 1. Regular Feature Release

```bash
git checkout main && git pull
make bump-minor
```

### 2. Critical Hotfix

```bash
# Fix the issue first
git add . && git commit -m "fix: critical bug"
make bump-patch
```

### 3. Pre-release Testing

```bash
# Alpha testing
git tag v2.0.0-alpha.1 && git push origin v2.0.0-alpha.1

# Beta testing
git tag v2.0.0-beta.1 && git push origin v2.0.0-beta.1

# Release candidate
git tag v2.0.0-rc.1 && git push origin v2.0.0-rc.1

# Final release
make bump-major
```

### 4. Development Cycle

```bash
# Feature development
git checkout -b feature/new-thing
git commit -m "feat: add new thing"

# Release when ready
git checkout main
git merge feature/new-thing
make bump-minor
```

## 📊 Version Types & Examples

| Version          | Type        | When to Use              |
| ---------------- | ----------- | ------------------------ |
| `v0.1.0`         | Development | Initial development      |
| `v1.0.0`         | Stable      | First production release |
| `v1.0.1`         | Patch       | Bug fixes only           |
| `v1.1.0`         | Minor       | New features added       |
| `v2.0.0`         | Major       | Breaking changes         |
| `v1.0.0-alpha.1` | Pre-release | Early testing            |
| `v1.0.0-beta.1`  | Pre-release | Feature complete testing |
| `v1.0.0-rc.1`    | Pre-release | Release candidate        |

## 🎯 Decision Guide

**When to bump PATCH (1.0.0 → 1.0.1):**

- 🐛 Bug fixes
- 🔒 Security patches
- 📝 Documentation fixes
- 🧪 Test improvements

**When to bump MINOR (1.0.0 → 1.1.0):**

- ✨ New features (backward compatible)
- 🗑️ Deprecating functionality
- 📦 New dependencies
- 🔧 Internal improvements

**When to bump MAJOR (1.0.0 → 2.0.0):**

- 💥 Breaking API changes
- 🗑️ Removing features
- 📝 Changing CLI interface
- 🔄 Major architecture changes

## 🛠️ Troubleshooting

### Fix Wrong Tag

```bash
git tag -d v1.0.0                    # Delete local
git push origin :refs/tags/v1.0.0   # Delete remote
git tag v1.0.1 && git push origin v1.0.1  # Recreate
```

### Clean State Issues

```bash
git status                           # Check what's dirty
git add . && git commit -m "msg"     # Commit changes
# or
git stash                            # Stash changes
```

### Version Detection Issues

```bash
git describe --tags --always --dirty # What Git sees
make version                         # What build sees
git tag -l                          # List all tags
```

## ✅ Release Checklist

- [ ] 🧪 **Tests pass**: `make test`
- [ ] 📝 **Docs updated**: CHANGELOG.md, README.md
- [ ] 🔍 **Clean state**: `git status` shows clean
- [ ] 🌿 **Correct branch**: Usually `main`
- [ ] 🏷️ **Right bump type**: patch/minor/major
- [ ] 🔨 **Build works**: `make build`
- [ ] ✔️ **Version correct**: `./bin/gaia-mcp-go --version`

## 💡 Pro Tips

### Commit Message Format

```bash
feat: new feature     # → MINOR bump
fix: bug fix         # → PATCH bump
feat!: breaking      # → MAJOR bump
docs: documentation  # → No bump
```

### Pre-release Strategy

```bash
# Progressive testing
v2.0.0-alpha.1  # Internal testing
v2.0.0-alpha.2  # Fix issues
v2.0.0-beta.1   # External testing
v2.0.0-rc.1     # Final validation
v2.0.0          # Release
```

### Version Planning

- 🏁 **v0.x.x**: Development/experimentation
- 🎯 **v1.0.0**: First stable production release
- 📈 **v1.x.x**: Stable evolution
- 🚀 **v2.0.0**: Major evolution

## 📚 More Details

For comprehensive explanations and examples, see [`semantic-versioning.md`](./semantic-versioning.md).
