#!/usr/bin/env bash
set -euo pipefail

OUT_ROOT=dist/licenses
OUT_THIRD="$OUT_ROOT/third_party_licenses"
MERGED_NOTICE="$OUT_ROOT/NOTICE"

rm -rf "$OUT_ROOT"
mkdir -p "$OUT_THIRD"

# 依存解決
go mod download

# go-licenses の導入
if ! command -v go-licenses >/dev/null 2>&1; then
  go install github.com/google/go-licenses@latest
fi

# 依存ライセンス全文を収集（各モジュールの LICENSE* / COPYING* 等）
go-licenses save ./... --save_path="$OUT_THIRD"

# 上流 NOTICE をマージ（存在するものだけ）
: > "$MERGED_NOTICE"
echo "Third-party notices for this distribution." >> "$MERGED_NOTICE"
echo >> "$MERGED_NOTICE"

# go list で各モジュールのディレクトリを取得
while IFS=$' ' read -r path ver dir; do
  # 置き換えや標準ライブラリなど dir が空のものはスキップ
  [[ -z "${dir}" || "${dir}" == " " ]] && continue
  
  # NOTICE / NOTICE.txt / NOTICE.md を探索
  for nf in "NOTICE" "NOTICE.txt" "Notice" "NOTICE.md"; do
    if [[ -f "${dir}/${nf}" ]]; then
      echo "----- NOTICE from ${path}@${ver} -----" >> "$MERGED_NOTICE"
      cat "${dir}/${nf}" >> "$MERGED_NOTICE"
      echo -e "\n" >> "$MERGED_NOTICE"
      break
    fi
  done
done < <(go list -m -f '{{.Path}} {{.Version}} {{.Dir}}' all)

# NOTICE が空なら削除（=依存に NOTICE が無いケース）
if [[ ! -s "$MERGED_NOTICE" || $(wc -c <"$MERGED_NOTICE") -lt 50 ]]; then
  rm -f "$MERGED_NOTICE"
fi

# 参考用レポート
go-licenses report ./... | tee "$OUT_ROOT/REPORT.txt"
