#!/usr/bin/env bash
set -euo pipefail

# rename.sh — Rename this template repo to your project name.
#
# Usage:
#   ./scripts/rename.sh <new-name> <github-org>
#
# Example:
#   ./scripts/rename.sh myproject myorg
#
# This replaces all occurrences of the current sentinel values with your chosen names.
# The script derives all case variants from the two arguments:
#
#   new-name  → myproject (slug-ns), my-project (kebab), MyProject (pascal),
#               My Project (human), MY_PROJECT (screaming), myproject.dev (domain),
#               com.myproject (bundle prefix)
#   github-org → used in Go module path: github.com/<org>/<name>
#
# The script is idempotent with respect to the *current* values in the repo.
# It detects what the current values are from backend/deploy/environments/prod/terraform/branding.tf.

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

# ---------------------------------------------------------------------------
# Parse arguments
# ---------------------------------------------------------------------------

if [[ $# -lt 2 ]]; then
  echo "Usage: $0 <new-name> <github-org>"
  echo ""
  echo "  <new-name>   Project name as a single lowercase word (e.g. 'myproject')"
  echo "  <github-org>  GitHub organization/user (e.g. 'myorg')"
  echo ""
  echo "All case variants (kebab, pascal, screaming, etc.) are derived automatically."
  exit 1
fi

NEW_SLUG_NS="$1"       # e.g. "zhuzh" or "myproject"
NEW_GITHUB_ORG="$2"    # e.g. "verygoodsoftwarenotvirus" or "myorg"

# ---------------------------------------------------------------------------
# Detect current values from branding.tf
# ---------------------------------------------------------------------------

BRANDING_TF="$REPO_ROOT/backend/deploy/environments/prod/terraform/branding.tf"
if [[ ! -f "$BRANDING_TF" ]]; then
  echo "Error: Cannot find $BRANDING_TF — are you in the right repo?"
  exit 1
fi

# Extract current values
OLD_COMPANY_NAME=$(grep 'company_name' "$BRANDING_TF" | head -1 | sed 's/.*= *"\(.*\)".*/\1/')
OLD_SLUG=$(grep 'company_slug ' "$BRANDING_TF" | head -1 | sed 's/.*= *"\(.*\)".*/\1/')
OLD_SLUG_NS=$(grep 'company_slug_ns' "$BRANDING_TF" | head -1 | sed 's/.*= *"\(.*\)".*/\1/')

# Detect current GitHub org from go.mod
GOMOD="$REPO_ROOT/backend/go.mod"
if [[ -f "$GOMOD" ]]; then
  OLD_GITHUB_ORG=$(head -1 "$GOMOD" | sed 's|module github.com/\([^/]*\)/.*|\1|')
  OLD_MODULE_REPO=$(head -1 "$GOMOD" | sed 's|module github.com/[^/]*/\([^/]*\)/.*|\1|')
else
  echo "Error: Cannot find $GOMOD"
  exit 1
fi

echo "=== Current values ==="
echo "  Company name:   $OLD_COMPANY_NAME"
echo "  Slug (kebab):   $OLD_SLUG"
echo "  Slug (compact): $OLD_SLUG_NS"
echo "  GitHub org:     $OLD_GITHUB_ORG"
echo "  Module repo:    $OLD_MODULE_REPO"
echo ""

# ---------------------------------------------------------------------------
# Derive all case variants for the new name
# ---------------------------------------------------------------------------

# Pascal case: capitalize first letter
NEW_PASCAL="$(echo "${NEW_SLUG_NS:0:1}" | tr '[:lower:]' '[:upper:]')${NEW_SLUG_NS:1}"

# Human-readable: same as pascal with spaces if the original had them,
# but since the new name is a single word, just use pascal
NEW_HUMAN="$NEW_PASCAL"

# Kebab case: for a single-word name, kebab = slug_ns
NEW_KEBAB="$NEW_SLUG_NS"

# Screaming snake: uppercase
NEW_SCREAMING="$(echo "$NEW_SLUG_NS" | tr '[:lower:]' '[:upper:]')"

# Domain
NEW_DOMAIN="${NEW_SLUG_NS}.dev"

# Derive old variants
OLD_PASCAL="$(echo "${OLD_SLUG_NS:0:1}" | tr '[:lower:]' '[:upper:]')${OLD_SLUG_NS:1}"
OLD_SCREAMING="$(echo "$OLD_SLUG_NS" | tr '[:lower:]' '[:upper:]' | tr '-' '_')"
# Handle the case where old screaming uses the kebab form with underscores
OLD_SCREAMING_FROM_KEBAB="$(echo "$OLD_SLUG" | tr '[:lower:]' '[:upper:]' | tr '-' '_')"
# Detect old domain TLD from branding.tf (e.g. ".dev" from "${local.company_slug_ns}.dev")
OLD_DOMAIN_TLD=$(grep 'public_domain' "$BRANDING_TF" | head -1 | sed 's/.*\.\([a-z]*\)".*/\1/')
OLD_DOMAIN="${OLD_SLUG_NS}.${OLD_DOMAIN_TLD:-dev}"

echo "=== New values ==="
echo "  Slug (compact): $NEW_SLUG_NS"
echo "  Slug (kebab):   $NEW_KEBAB"
echo "  Pascal:         $NEW_PASCAL"
echo "  Human:          $NEW_HUMAN"
echo "  Screaming:      $NEW_SCREAMING"
echo "  Domain:         $NEW_DOMAIN"
echo "  GitHub org:     $NEW_GITHUB_ORG"
echo ""

# ---------------------------------------------------------------------------
# Perform replacements
# ---------------------------------------------------------------------------

# We use find + sed. On macOS, sed -i requires '' as backup extension.
# On Linux, sed -i works without it. Detect platform.
if [[ "$(uname)" == "Darwin" ]]; then
  SED_INPLACE=(sed -i '')
else
  SED_INPLACE=(sed -i)
fi

# Find all text files, excluding .git, node_modules, vendor, and binary files
find_text_files() {
  find "$REPO_ROOT" \
    -not -path '*/.git/*' \
    -not -path '*/node_modules/*' \
    -not -path '*/vendor/*' \
    -not -path '*/artifacts/*' \
    -not -path '*/.terraform/*' \
    -not -name '*.png' \
    -not -name '*.jpg' \
    -not -name '*.jpeg' \
    -not -name '*.gif' \
    -not -name '*.ico' \
    -not -name '*.woff' \
    -not -name '*.woff2' \
    -not -name '*.ttf' \
    -not -name '*.eot' \
    -not -name '*.zip' \
    -not -name '*.tar.gz' \
    -not -name '*.jar' \
    -not -name '*.lock' \
    -not -name 'package-lock.json' \
    -not -name 'go.sum' \
    -type f
}

replace_in_files() {
  local old="$1"
  local new="$2"
  local desc="$3"

  if [[ "$old" == "$new" ]]; then
    echo "  SKIP: $desc (already matches)"
    return
  fi

  # Use grep to find files containing the pattern, then sed to replace
  local count
  count=$(grep -rl --fixed-strings "$old" "$REPO_ROOT" \
    --exclude-dir=.git \
    --exclude-dir=node_modules \
    --exclude-dir=vendor \
    --exclude-dir=artifacts \
    --exclude-dir=.terraform \
    --exclude='*.png' \
    --exclude='*.jpg' \
    --exclude='*.ico' \
    --exclude='*.woff*' \
    --exclude='package-lock.json' \
    --exclude='go.sum' \
    2>/dev/null | wc -l | tr -d ' ')

  if [[ "$count" -eq 0 ]]; then
    echo "  SKIP: $desc (no matches)"
    return
  fi

  echo "  $desc: $old → $new ($count files)"

  grep -rl --fixed-strings "$old" "$REPO_ROOT" \
    --exclude-dir=.git \
    --exclude-dir=node_modules \
    --exclude-dir=vendor \
    --exclude-dir=artifacts \
    --exclude-dir=.terraform \
    --exclude='*.png' \
    --exclude='*.jpg' \
    --exclude='*.ico' \
    --exclude='*.woff*' \
    --exclude='package-lock.json' \
    --exclude='go.sum' \
    2>/dev/null | while IFS= read -r file; do
      # Use perl for reliable fixed-string replacement with | delimiter to avoid / conflicts
      perl -pi -e "s|\Q${old}\E|${new}|g" "$file"
    done
}

echo "=== Replacing patterns (most specific first) ==="

# 1. Go module path (most specific — contains org + repo)
replace_in_files \
  "github.com/${OLD_GITHUB_ORG}/${OLD_MODULE_REPO}" \
  "github.com/${NEW_GITHUB_ORG}/${NEW_SLUG_NS}" \
  "Go module path"

# 2. npm package scope
replace_in_files \
  "@${OLD_SLUG_NS}/" \
  "@${NEW_SLUG_NS}/" \
  "npm package scope"

# 3. Screaming snake (env var prefix) — from kebab form (ZHUZH)
replace_in_files \
  "${OLD_SCREAMING_FROM_KEBAB}" \
  "${NEW_SCREAMING}" \
  "Screaming snake (env vars)"

# 4. Human-readable name
replace_in_files \
  "${OLD_COMPANY_NAME}" \
  "${NEW_HUMAN}" \
  "Human-readable name"

# 5. PascalCase branding
#    Handle the terraform branding.tf special case where it uses PascalCase without spaces
OLD_PASCAL_NO_SPACE="$(echo "$OLD_SLUG" | sed 's/-/ /g' | awk '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) tolower(substr($i,2))}1' | tr -d ' ')"
if [[ "$OLD_PASCAL_NO_SPACE" != "$OLD_COMPANY_NAME" ]]; then
  replace_in_files \
    "${OLD_PASCAL_NO_SPACE}" \
    "${NEW_PASCAL}" \
    "PascalCase name"
fi

# 6. Domain name (zhuzh.dev → zhuzh.dev)
replace_in_files \
  "${OLD_DOMAIN}" \
  "${NEW_DOMAIN}" \
  "Domain name"

# 7. iOS bundle ID prefix (com.zhuzh → com.zhuzh)
replace_in_files \
  "com.${OLD_SLUG_NS}" \
  "com.${NEW_SLUG_NS}" \
  "iOS bundle ID prefix"

# 8. Kebab-case slug (zhuzh → zhuzh)
replace_in_files \
  "${OLD_SLUG}" \
  "${NEW_KEBAB}" \
  "Kebab-case slug"

# 9. Compact slug catch-all (zhuzh → zhuzh)
#    This must be LAST to avoid clobbering more-specific patterns
replace_in_files \
  "${OLD_SLUG_NS}" \
  "${NEW_SLUG_NS}" \
  "Compact slug (catch-all)"

echo ""
echo "=== Post-rename steps ==="
echo "  1. cd backend && go mod tidy"
echo "  2. cd backend && go build ./..."
echo "  3. cd frontend && npm install"
echo "  4. Regenerate proto code if proto tooling is available"
echo "  5. Regenerate sqlc code: cd backend && make querier"
echo "  6. Regenerate env vars: cd backend && go run ./cmd/tools/codegen/valid_env_vars"
echo ""
echo "Done!"
