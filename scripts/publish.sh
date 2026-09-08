#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
publish="$PWD/.cache/publish"
if git ls-remote --exit-code --heads origin srs > /dev/null; then
  git fetch origin srs
  git worktree add --detach "$publish" FETCH_HEAD
else
  git worktree add --detach "$publish" HEAD
  git -C "$publish" checkout --orphan srs
  git -C "$publish" rm -rf .
fi
rm -rf "$publish/geoip" "$publish/geosite"
cp -R .cache/output/. "$publish/"
git -C "$publish" add --all
if git -C "$publish" diff --cached --quiet; then
  echo 'No changes; preserving current srs commit.'
  exit 0
fi
git -C "$publish" -c user.name='github-actions[bot]' \
  -c user.email='41898282+github-actions[bot]@users.noreply.github.com' \
  commit -m 'Update SRS rules from verified upstream release'
git -C "$publish" push origin HEAD:refs/heads/srs
