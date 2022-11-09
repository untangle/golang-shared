import subprocess
import re
import sys

tag_reg = r'^v([0-9]+)\.([0-9]+)\.([0-9]+)$'
patterns = {
    'major': (r'^version: major\s*$', 0),
    'minor': (r'^version: minor\s*$', 1),
    'bug': (r'^version: bug\s*$', 2),
}


def tag2version(tag: str):
    output = re.match(tag_reg, tag)
    if output is not None:
        return (int(output.group(1)),
                int(output.group(2)),
                int(output.group(3)))
    return None


def find_latest_tag():
    subprocess.run(['git', 'fetch', 'origin', '--tags'], check=True)
    tags = [tag2version(i) for i in
            (subprocess.run(['git', 'tag'], stdout=subprocess.PIPE,
                            check=True)
             .stdout.decode().splitlines())]
    latest = sorted([i for i in tags if i])[-1]
    return latest


latest = find_latest_tag()
msg = sys.stdin.read().splitlines()
for type, (pattern, index) in patterns.items():
    if any((re.match(pattern, i) for i in msg)):
        new = list(latest)
        new[index] += 1
        for i in range(index+1, len(new)):
            new[i] = 0
        print(f'bumping with a {type} version to: {new}', file=sys.stderr)
        print('v' + '.'.join([str(i) for i in new]))
        break
else:
    print(f'Could not read commit msg: {msg}!', file=sys.stderr)
    raise SystemExit(1)
