#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
fasd_path="${FASDER_BENCH_FASD:-"$HOME/workspace/clvv-fasd/fasd"}"
entries="${FASDER_BENCH_ENTRIES:-2000}"
count="${FASDER_BENCH_COUNT:-5}"

if [[ ! -x "$fasd_path" ]]; then
  echo "Original fasd script not found or not executable: $fasd_path" >&2
  echo "Set FASDER_BENCH_FASD=/path/to/fasd and retry." >&2
  exit 1
fi

tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

fasder_bin="$tmpdir/fasder"
gocache="${GOCACHE:-"$tmpdir/gocache"}"

cd "$repo_root"
GOCACHE="$gocache" go build -buildvcs=false -o "$fasder_bin" .

FASDER_BENCH_FASDER="$fasder_bin" \
FASDER_BENCH_FASD="$fasd_path" \
FASDER_BENCH_ENTRIES="$entries" \
GOCACHE="$gocache" \
go test -buildvcs=false -run '^$' -bench '^BenchmarkCliCompare' -count "$count" "$@"
