# Release a new version
release version:
    #!/usr/bin/env bash
    set -euo pipefail

    echo "Creating release {{version}}..."

    # Ensure we're on main and up to date
    git checkout main
    git pull origin main

    # Commit any changes
    if [[ -n $(git status -s) ]]; then
        git add .
        git commit -m "Prepare release {{version}}"
        git push origin main
    fi

    # Delete tag if it exists locally
    git tag -d {{version}} 2>/dev/null || true

    # Delete tag if it exists remotely
    git push origin :refs/tags/{{version}} 2>/dev/null || true

    # Create and push new tag
    git tag -a {{version}} -m "Release {{version}}"
    git push origin {{version}}

    # Run goreleaser
    goreleaser release --clean

    echo "✅ Release {{version}} completed!"

# Release a new beta version (auto-increments)
release-beta:
    #!/usr/bin/env bash
    set -euo pipefail

    # Get latest beta tag or default to beta.0
    LATEST=$(git tag -l "v*-beta.*" | sort -V | tail -n1)

    if [[ -z "$LATEST" ]]; then
        NEXT="v1.0.0-beta.1"
    else
        # Extract version parts
        VERSION=$(echo $LATEST | sed 's/v\(.*\)-beta\.\(.*\)/\1/')
        BETA_NUM=$(echo $LATEST | sed 's/v\(.*\)-beta\.\(.*\)/\2/')
        NEXT_BETA=$((BETA_NUM + 1))
        NEXT="v${VERSION}-beta.${NEXT_BETA}"
    fi

    echo "Next beta version: $NEXT"
    just release $NEXT

# Release a stable version
release-stable version:
    just release v{{version}}

# List all releases
list-releases:
    git tag -l "v*" | sort -V

# Clean up failed release
cleanup-release version:
    #!/usr/bin/env bash
    set -euo pipefail

    echo "Cleaning up {{version}}..."
    git tag -d {{version}} 2>/dev/null || true
    git push origin :refs/tags/{{version}} 2>/dev/null || true
    echo "✅ Cleanup completed"
