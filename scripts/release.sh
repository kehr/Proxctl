#!/usr/bin/env bash
set -euo pipefail

remote="${REMOTE:-origin}"
branch="${BRANCH:-main}"
bump="${BUMP:-patch}"
explicit_version="${RELEASE_VERSION:-${VERSION:-}}"
dry_run="${DRY_RUN:-0}"

usage() {
  cat <<'EOF'
Usage:
  make release                         # bump patch, commit, tag, push
  make release BUMP=minor              # bump minor
  make release BUMP=major              # bump major
  make release RELEASE_VERSION=v1.2.3  # use explicit version
  make release-dry-run [BUMP=patch]

Environment:
  BUMP=patch|minor|major
  RELEASE_VERSION=vX.Y.Z
  REMOTE=origin
  BRANCH=main
  DRY_RUN=1
EOF
}

run() {
  if [[ "$dry_run" == "1" ]]; then
    printf 'dry-run: %q' "$1"
    shift
    for arg in "$@"; do
      printf ' %q' "$arg"
    done
    printf '\n'
  else
    "$@"
  fi
}

fail() {
  echo "release: $*" >&2
  exit 1
}

current_version="$(tr -d '[:space:]' < VERSION)"

next_version() {
  local version="$1"
  local part="$2"
  [[ "$version" =~ ^v([0-9]+)\.([0-9]+)\.([0-9]+)$ ]] || fail "VERSION must be vMAJOR.MINOR.PATCH, got $version"
  local major="${BASH_REMATCH[1]}"
  local minor="${BASH_REMATCH[2]}"
  local patch="${BASH_REMATCH[3]}"
  case "$part" in
    patch) patch=$((patch + 1)) ;;
    minor) minor=$((minor + 1)); patch=0 ;;
    major) major=$((major + 1)); minor=0; patch=0 ;;
    *) fail "BUMP must be patch, minor, or major" ;;
  esac
  printf 'v%s.%s.%s\n' "$major" "$minor" "$patch"
}

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
  usage
  exit 0
fi

if [[ -n "$explicit_version" ]]; then
  target_version="$explicit_version"
  [[ "$target_version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]] || fail "RELEASE_VERSION must be vMAJOR.MINOR.PATCH"
else
  target_version="$(next_version "$current_version" "$bump")"
fi

current_branch="$(git branch --show-current)"
[[ "$current_branch" == "$branch" ]] || fail "must release from $branch, current branch is $current_branch"

if [[ -n "$(git status --porcelain)" ]]; then
  fail "working tree is not clean; commit or stash changes first"
fi

git fetch "$remote" --tags
counts="$(git rev-list --left-right --count "$remote/$branch"...HEAD)"
set -- $counts
behind="${1:-0}"
if (( behind > 0 )); then
  fail "local $branch is behind $remote/$branch; pull before release"
fi

if git rev-parse -q --verify "refs/tags/$target_version" >/dev/null; then
  fail "local tag already exists: $target_version"
fi
if git ls-remote --exit-code --tags "$remote" "refs/tags/$target_version" >/dev/null 2>&1; then
  fail "remote tag already exists: $target_version"
fi

echo "release: $current_version -> $target_version"

if [[ "$dry_run" == "1" ]]; then
  cat <<EOF
dry-run: would write VERSION=$target_version
dry-run: would run make verify
dry-run: would commit VERSION with message: chore: release $target_version
dry-run: would create annotated tag: $target_version
dry-run: would push $remote $branch
dry-run: would push $remote $target_version
EOF
  exit 0
fi

echo "$target_version" > VERSION

make verify

run git add VERSION
run git commit -m "chore: release $target_version"
run git tag -a "$target_version" -m "Proxctl $target_version"
run git push "$remote" "$branch"
run git push "$remote" "$target_version"

cat <<EOF
release: pushed $target_version
release: GitHub Actions will build and publish the release from the tag.
release: https://github.com/kehr/Proxctl/releases/tag/$target_version
EOF
