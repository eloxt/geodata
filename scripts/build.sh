#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
root="$PWD"
converter=368e4b324bb4b8c90da93d1b8429bd0af337dfda
mkdir -p .cache/bin
rm -rf .cache/input .cache/output
mkdir -p .cache/input .cache/output/geosite .cache/output/geoip

# Resolve one immutable release tag so DAT and checksum downloads cannot mix releases.
gh api repos/Loyalsoldier/v2ray-rules-dat/releases/latest > .cache/release.json
tag=$(python3 -c 'import json; print(json.load(open(".cache/release.json"))["tag_name"])')
gh release download "$tag" --repo Loyalsoldier/v2ray-rules-dat \
  --pattern '*.dat' --pattern '*.dat.sha256sum' --dir .cache/input
(cd .cache/input && sha256sum --check geoip.dat.sha256sum geosite.dat.sha256sum)

GOBIN="$root/.cache/bin" go install "github.com/runetfreedom/geodat2srs@$converter"
.cache/bin/geodat2srs geoip -i .cache/input/geoip.dat -o .cache/output/geoip --prefix ''
.cache/bin/geodat2srs geosite -i .cache/input/geosite.dat -o .cache/output/geosite --prefix ''
# The converter exports all categories; publish only the two CN rule sets.
find .cache/output -type f -name '*.srs' ! -name 'cn.srs' -delete
go run ./scripts/verify.go .cache/input .cache/output
python3 scripts/manifest.py
