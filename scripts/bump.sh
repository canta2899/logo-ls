#!/usr/bin/env bash
set -euo pipefail

BUMP_TYPE="${1:-}"

if [[ "$BUMP_TYPE" != "major" && "$BUMP_TYPE" != "minor" && "$BUMP_TYPE" != "patch" ]]; then
  echo "Usage: $0 [major|minor|patch]" >&2
  exit 1
fi

if ! git diff --quiet || ! git diff --cached --quiet; then
  echo "Error: there are uncommitted changes. Please commit or stash them first." >&2
  exit 1
fi

LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
VERSION="${LATEST_TAG#v}"

IFS='.' read -r MAJOR MINOR PATCH <<<"$VERSION"

case "$BUMP_TYPE" in
major)
  MAJOR=$((MAJOR + 1))
  MINOR=0
  PATCH=0
  ;;
minor)
  MINOR=$((MINOR + 1))
  PATCH=0
  ;;
patch) PATCH=$((PATCH + 1)) ;;
esac

NEW_VERSION="v${MAJOR}.${MINOR}.${PATCH}"
CHANGELOG="CHANGELOG.md"

if ! grep -qi '## \[Unreleas' "$CHANGELOG"; then
  echo "Error: No [Unreleased] section found in $CHANGELOG" >&2
  exit 1
fi

# Replace the [Unreleased] heading with versioned heading and prepend a fresh [Unreleased] block
perl -i -pe "s/## \[Unreleas[^\]]*\]/## [Unreleased]\n\n## logo-ls [${MAJOR}.${MINOR}.${PATCH}]/i" "$CHANGELOG"

echo "Updated $CHANGELOG: [Unreleased] -> logo-ls [${MAJOR}.${MINOR}.${PATCH}]"

git add "$CHANGELOG"
git commit -m "chore: bumped to ${NEW_VERSION}"
git tag "${NEW_VERSION}"

echo "\nTagged ${NEW_VERSION}. Run the following to release:\n"
echo "  git push --follow-tags"
