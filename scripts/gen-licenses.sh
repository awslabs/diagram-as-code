#!/usr/bin/env bash
set -euo pipefail

OUT_ROOT="${1:-dist/licenses}"
OUT_THIRD="$OUT_ROOT/third_party_licenses"
MERGED_NOTICE="$OUT_ROOT/NOTICE"

rm -rf "$OUT_ROOT"
mkdir -p "$OUT_ROOT"

# Download dependencies
go mod download

# Install go-licenses if not available
if ! command -v go-licenses >/dev/null 2>&1; then
  go install github.com/google/go-licenses@latest
fi

# Collect full license texts for dependencies (LICENSE*, COPYING*, etc.)
go-licenses save ./... --save_path="$OUT_THIRD" --ignore="github.com/golang/freetype"

# Merge upstream NOTICE files (only existing ones)
: > "$MERGED_NOTICE"
echo "Third-party notices for this distribution." >> "$MERGED_NOTICE"
echo >> "$MERGED_NOTICE"

# Get directory for each module using go list
while IFS=$' ' read -r path ver dir; do
  # Skip modules with empty dir (replacements, standard library, etc.)
  [[ -z "${dir}" || "${dir}" == " " ]] && continue
  
  # Search for NOTICE / NOTICE.txt / NOTICE.md files
  for nf in "NOTICE" "NOTICE.txt" "Notice" "NOTICE.md"; do
    if [[ -f "${dir}/${nf}" ]]; then
      echo "----- NOTICE from ${path}@${ver} -----" >> "$MERGED_NOTICE"
      cat "${dir}/${nf}" >> "$MERGED_NOTICE"
      echo -e "\n" >> "$MERGED_NOTICE"
      break
    fi
  done
done < <(go list -m -f '{{.Path}} {{.Version}} {{.Dir}}' all)

# Remove NOTICE if empty (no dependencies have NOTICE files)
if [[ ! -s "$MERGED_NOTICE" || $(wc -c <"$MERGED_NOTICE") -lt 50 ]]; then
  rm -f "$MERGED_NOTICE"
fi

# Generate reference report
go-licenses report ./... | tee "$OUT_ROOT/REPORT.txt"
