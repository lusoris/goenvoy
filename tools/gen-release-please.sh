#!/usr/bin/env bash
# Regenerate release-please-config.json and .release-please-manifest.json
# from the current set of go.mod files + existing git tags.
#
# Run from repo root: tools/gen-release-please.sh
set -euo pipefail

CONFIG=release-please-config.json
MANIFEST=.release-please-manifest.json

mapfile -t MODS < <(find . -name 'go.mod' -not -path './.workingdir*/*' -exec dirname {} \; | sed 's|^\./||' | sort)

# --- config ---
{
  cat <<'EOF'
{
  "$schema": "https://raw.githubusercontent.com/googleapis/release-please/main/schemas/config.json",
  "release-type": "go",
  "bump-minor-pre-major": true,
  "bump-patch-for-minor-pre-major": false,
  "include-v-in-tag": true,
  "separate-pull-requests": true,
  "tag-separator": "/",
  "changelog-sections": [
    { "type": "feat", "section": "Features" },
    { "type": "fix", "section": "Bug Fixes" },
    { "type": "perf", "section": "Performance" },
    { "type": "refactor", "section": "Code Refactoring" },
    { "type": "deps", "section": "Dependencies" },
    { "type": "docs", "section": "Documentation", "hidden": true },
    { "type": "test", "section": "Tests", "hidden": true },
    { "type": "ci", "section": "CI", "hidden": true },
    { "type": "build", "section": "Build", "hidden": true },
    { "type": "chore", "section": "Chores", "hidden": true },
    { "type": "style", "section": "Style", "hidden": true },
    { "type": "revert", "section": "Reverts", "hidden": true }
  ],
  "packages": {
EOF

  last=$((${#MODS[@]} - 1))
  for i in "${!MODS[@]}"; do
    mod="${MODS[$i]}"
    sep=","
    [ "$i" -eq "$last" ] && sep=""
    printf '    "%s": { "package-name": "goenvoy/%s", "component": "%s" }%s\n' \
      "$mod" "$mod" "$mod" "$sep"
  done

  cat <<'EOF'
  }
}
EOF
} > "$CONFIG"

# --- manifest (latest tag per module, stripped of leading 'v') ---
{
  echo '{'
  last=$((${#MODS[@]} - 1))
  for i in "${!MODS[@]}"; do
    mod="${MODS[$i]}"
    sep=","
    [ "$i" -eq "$last" ] && sep=""
    latest=$(git tag --list "${mod}/v*" --sort=-v:refname | head -1 || true)
    if [ -n "$latest" ]; then
      ver="${latest##*/v}"
    else
      ver="0.0.0"
    fi
    printf '  "%s": "%s"%s\n' "$mod" "$ver" "$sep"
  done
  echo '}'
} > "$MANIFEST"

echo "wrote $CONFIG and $MANIFEST (${#MODS[@]} packages)"
