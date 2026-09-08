import hashlib
import json
import os
from pathlib import Path

release = json.loads(Path('.cache/release.json').read_text())
output = Path('.cache/output')
files = {}
for path in sorted(output.glob('*/*.srs')):
    files[path.relative_to(output).as_posix()] = hashlib.sha256(path.read_bytes()).hexdigest()
manifest = {
    'upstream': 'Loyalsoldier/v2ray-rules-dat',
    'release': release['tag_name'],
    'release_url': release['html_url'],
    'converter': 'runetfreedom/geodat2srs@368e4b324bb4b8c90da93d1b8429bd0af337dfda',
    'sources': {
        name: hashlib.sha256((Path('.cache/input') / name).read_bytes()).hexdigest()
        for name in ('geoip.dat', 'geosite.dat')
    },
    'counts': {kind: sum(name.startswith(kind + '/') for name in files) for kind in ('geoip', 'geosite')},
    'files': files,
}
(output / 'manifest.json').write_text(json.dumps(manifest, ensure_ascii=False, indent=2) + '\n')
(output / 'SHA256SUMS').write_text(''.join(f'{digest}  {name}\n' for name, digest in files.items()))
(output / 'README.md').write_text(
    '# sing-box SRS rules\n\n'
    'Generated from [Loyalsoldier/v2ray-rules-dat](' + release['html_url'] + ').\n\n'
    f'Upstream release: `{release["tag_name"]}`.\n\n'
    'Only CN rules are published: `geosite/cn.srs` and `geoip/cn.srs`.\n\n'
    'See [source and usage](https://github.com/eloxt/geodata) and `manifest.json` for provenance and checksums. '
    'Upstream data retains its original licenses.\n'
)
print(json.dumps(manifest['counts']))
if os.environ.get('GITHUB_STEP_SUMMARY'):
    with open(os.environ['GITHUB_STEP_SUMMARY'], 'a') as summary:
        summary.write(f'Upstream: {release["html_url"]}\n\nConverted: {manifest["counts"]}\n')
