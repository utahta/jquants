#!/usr/bin/env bash
# J-Quants 公式ドキュメント（Markdown配信）を docs/v2/ に同期する。
# llms.txt に載っているページをシードに、ページ内の /ja/spec/ リンクも辿って取得する。
# docs/v2/ はローカルキャッシュであり、リポジトリにはコミットしない（.gitignore対象）。
set -euo pipefail

BASE="https://jpx-jquants.com"
DEST="$(cd "$(dirname "$0")/.." && pwd)/docs/v2"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

queue="$(curl -fsS "$BASE/llms.txt" | grep -oE '(https://jpx-jquants\.com)?/ja/spec/[a-z0-9/_-]+' | sed "s|^$BASE||" | sort -u)"
seen=""

while [ -n "$queue" ]; do
  next=""
  for path in $queue; do
    case " $seen " in *" $path "*) continue ;; esac
    seen="$seen $path"
    out="$TMP/${path#/ja/spec/}.md"
    mkdir -p "$(dirname "$out")"
    if ! curl -fsS "$BASE$path.md" -o "$out"; then
      echo "skip (fetch failed): $path" >&2
      rm -f "$out"
      continue
    fi
    echo "fetched: ${path#/ja/spec/}.md"
    links="$(grep -oE '\]\((https://jpx-jquants\.com)?/ja/spec/[a-z0-9/_-]+' "$out" || true)"
    next="$next $(printf '%s\n' "$links" | sed -e 's|^](||' -e "s|^$BASE||" | sort -u | tr '\n' ' ')"
  done
  queue="$(printf '%s\n' $next | sort -u)"
done

rm -rf "$DEST"
mkdir -p "$(dirname "$DEST")"
mv "$TMP" "$DEST"
trap - EXIT
echo "synced: $(find "$DEST" -name '*.md' | wc -l | tr -d ' ') files -> $DEST"
