#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get current version from git tags
get_current_version() {
    git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//' || echo "0.0.0"
}

# Parse version into components
parse_version() {
    local version=$1
    echo $version | sed 's/\([0-9]*\)\.\([0-9]*\)\.\([0-9]*\).*/\1 \2 \3/'
}

# Bump version based on type
bump_version() {
    local current_version=$1
    local bump_type=$2
    
    read major minor patch <<< $(parse_version $current_version)
    
    case $bump_type in
        major)
            major=$((major + 1))
            minor=0
            patch=0
            ;;
        minor)
            minor=$((minor + 1))
            patch=0
            ;;
        patch)
            patch=$((patch + 1))
            ;;
        *)
            echo -e "${RED}Error: Invalid bump type. Use 'major', 'minor', or 'patch'${NC}"
            exit 1
            ;;
    esac
    
    echo "$major.$minor.$patch"
}

# Main function
main() {
    local bump_type=${1:-patch}
    local dry_run=${2:-false}
    
    # Ensure we're in a git repository
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        echo -e "${RED}Error: Not in a git repository${NC}"
        exit 1
    fi
    
    # Ensure working directory is clean
    if [[ -n $(git status --porcelain) ]]; then
        echo -e "${RED}Error: Working directory not clean. Please commit your changes.${NC}"
        exit 1
    fi
    
    # Get current version
    current_version=$(get_current_version)
    echo -e "${YELLOW}Current version: v$current_version${NC}"
    
    # Calculate new version
    new_version=$(bump_version $current_version $bump_type)
    echo -e "${GREEN}New version: v$new_version${NC}"
    
    if [[ $dry_run == "true" ]]; then
        echo -e "${YELLOW}Dry run mode - no changes made${NC}"
        exit 0
    fi
    
    # Confirm with user
    read -p "Create tag v$new_version? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Cancelled."
        exit 0
    fi
    
    # Create and push tag
    git tag -a "v$new_version" -m "Release v$new_version"
    git push origin "v$new_version"
    
    echo -e "${GREEN}Successfully created and pushed tag v$new_version${NC}"
}

# Show usage if no arguments or help requested
if [[ $# -eq 0 ]] || [[ $1 == "-h" ]] || [[ $1 == "--help" ]]; then
    echo "Usage: $0 <bump_type> [dry_run]"
    echo ""
    echo "bump_type:"
    echo "  major  - Bump major version (breaking changes)"
    echo "  minor  - Bump minor version (new features)"
    echo "  patch  - Bump patch version (bug fixes)"
    echo ""
    echo "Options:"
    echo "  dry_run - Show what would happen without making changes"
    echo ""
    echo "Examples:"
    echo "  $0 patch              # Bump patch version"
    echo "  $0 minor              # Bump minor version"
    echo "  $0 major              # Bump major version"
    echo "  $0 patch dry_run      # Dry run for patch bump"
    exit 0
fi

main "$@"